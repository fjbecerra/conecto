package oauth

import (
	"conecto/auth/connections"
	"conecto/auth/credentials"
	"conecto/auth/oauth/state"
	"conecto/connectors"
	"conecto/sync"
	"context"
	"errors"
	"net/url"
)

type Service struct {
	connectionStore 	connections.Store
	credentialService 	credentials.CredentialService
	stateSigner	    	state.StateSigner      
	registry        	*connectors.Registry
	syncService 		*sync.SyncService
}

func NewService(
	connectionStore connections.Store,
	credentialService credentials.CredentialService,
	stateSigner state.StateSigner,
	registry *connectors.Registry,
	syncService *sync.SyncService) *Service {
		return &Service{
			connectionStore: connectionStore,
			credentialService: credentialService,
			stateSigner: stateSigner,
			registry: registry,
			syncService: syncService,
		}
}

func (s *Service) BeginAuthorization(ctx context.Context, connectionID string) (string, error) {

	connection, err := s.connectionStore.Get(ctx, connectionID)
	if err != nil {
		return "", err
	}

	connector,_ := s.registry.Get(connection.Provider)

	state, err := s.stateSigner.Sign(connectionID)

	if err != nil {
		return "", err
	}

	return connector.AuthorizeURL(
		ctx,
		connection,
		state,
	)
}

func (s *Service) HandleCallback(ctx context.Context, values url.Values) error {

	code := values.Get("code")
	if code == "" {
		return errors.New("missing oauth code")
	}

	state := values.Get("state")
	if state == "" {
		return errors.New("missing oauth state")
	}


	// 1. Verify state and recover connection ID
	connectionID, err := s.stateSigner.Verify(state)
	if err != nil {
		return err
	}


	// 2. Load connection
	connection, err :=s.connectionStore.Get(ctx,connectionID)

	if err != nil {
		return err
	}

	// 3. Get connector
	connector, err := s.registry.Get(connection.Provider)

	if err != nil {
		return errors.New("connector not found")
	}


	// 4. Exchange OAuth code for credentials
	credential, err := connector.Exchange(ctx,connection,code)

	if err != nil {
		return err
	}


	// 5. Save credential securely
	err = s.credentialService.Save(
			ctx,
			connection,
			credential,
	)

	//6. Update connection status
	s.connectionStore.UpdateStatus(ctx, connectionID, connections.StatusConnected )


	if err != nil {
		return err
	}

	//7. Start sync job - this should be a backfill of last 90 days
	return s.syncService.ScheduleConnectionSync(ctx,connection)
}
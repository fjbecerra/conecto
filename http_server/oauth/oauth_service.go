package oauth

import (
	"conecto/core"
	"conecto/http_server/oauth/state"
	"conecto/resources"
	"conecto/resources/base/api"
	"conecto/stores/connections"
	"conecto/stores/credentials"
	"conecto/sync"
	"context"
	"errors"
	"net/url"
)

type Service struct {
	connectionStore 	connections.Store
	credentialService 	credentials.CredentialService
	stateSigner	    	state.StateSigner      
	resourceRegistry  	*resources.ResourcesRegistry
	syncService 		*sync.SyncService
}

func NewService(
	connectionStore connections.Store,
	credentialService credentials.CredentialService,
	stateSigner state.StateSigner,
	resourceRegistry   *resources.ResourcesRegistry,
	syncService *sync.SyncService) *Service {
		return &Service{
			connectionStore: connectionStore,
			credentialService: credentialService,
			stateSigner: stateSigner,
			resourceRegistry: resourceRegistry,
			syncService: syncService,
		}
}

func (s *Service) BeginAuthorization(ctx context.Context, connectionID string) (string, error) {

	connection, err := s.connectionStore.Get(ctx, connectionID)
	if err != nil {
		return "", err
	}

	oauthentication:= s.resourceRegistry.Get(resources.ResourceName(connection.ResourceName))
	oauthProvider, ok := oauthentication.(api.OAuthProvider)
	if !ok {
		return "", errors.New("resource does not support OAuth")
	}

	state, err := s.stateSigner.Sign(connectionID)

	if err != nil {
		return "", err
	}

	return oauthProvider.AuthorizeURL(
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
	oauthentication:= s.resourceRegistry.Get(resources.ResourceName(connection.ResourceName))
	oauthProvider, ok := oauthentication.(api.OAuthProvider)
	if !ok {
		return errors.New("resource does not support OAuth")
	}
	

	// 4. Exchange OAuth code for credentials
	credential, err := oauthProvider.Exchange(ctx,connection,code)

	if err != nil {
		return err
	}


	// 5. Save credential securely
	err = s.credentialService.Save(
			ctx,
			connection,
			credential,
	)

	if err != nil {
		return err
	}

	//6. Update connection status
	s.connectionStore.UpdateStatus(ctx, connectionID, core.StatusConnected)


	if err != nil {
		return err
	}

	//7. Start sync job - this should be a backfill of last 90 days
	return s.syncService.ScheduleConnectionSync(ctx,connection)
}
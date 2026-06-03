package statestores

import "fmt"

type Status int

const (
    Running Status = iota
    Completed
    Failed
	Stopped
)

func (s Status) String() string {
    switch s {
    case Running:
        return "RUNNING"
    case Completed:
        return "COMPLETED"
	case Failed:
        return "FAILED"
	case Stopped:
        return "STOPPED"
    default:
        return "UNKNOWN"
    }
}

func ParseStatus(s string) (Status, error) {
    switch s {
    case "RUNNING":
        return Running, nil
    case "COMPLETED":
        return Completed, nil
	case "FAILED":
        return Failed, nil
	case "STOPPED":
        return Stopped, nil
    default:
        return 0, fmt.Errorf("unknown status: %s", s)
    }
}

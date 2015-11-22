package heartbeats

import (
	"fmt"
	"os"

	"github.com/azer/logger"
	"github.com/hackliff/serf/command/agent"
	"github.com/hackliff/serf/serf"
	"github.com/samalba/dockerclient"
)

const DEFAULT_DOCKER_HOST string = "unix:///var/run/docker.sock"

func init() {
	log.Info("registering heartbeat: docker")
	Add("docker", func(conf map[string]string) (Plugin, error) {
		host, ok := conf["host"]
		if !ok {
			host = os.Getenv("DOCKER_HOST")
			if host == "" {
				log.Info("no host information found, fallback to default: %s")
				host = DEFAULT_DOCKER_HOST
			}
		}

		log.Info("connecting to docker (%s)", host)
		// TODO support for tls
		docker, err := dockerclient.NewDockerClient(host, nil)
		if err != nil {
			return nil, err
		}

		return &DockerMonitor{
			docker: docker,
			logger: logger.New("sentinel.plugins.heartbeats.docker"),
		}, nil
	})
}

type DockerMonitor struct {
	docker *dockerclient.DockerClient
	logger *logger.Logger
}

// Callback used to listen to Docker's events
func (d DockerMonitor) eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	log.Info("Received event: %#v (%#v)", *event, args)
	e := serf.UserEvent{
		Name:     "heartbeat-docker",
		Payload:  []byte(fmt.Sprintf("")),
		Coalesce: true,
	}
	args[1].(agent.EventHandler).HandleEvent(e)
}

func (d DockerMonitor) Schedule(name string, handler agent.EventHandler) {
	log.Info("monitoring docker events: %s", name)
	version, err := d.docker.Version()
	if err != nil {
		log.Error("error gathering docker info: %s", err)
	} else {
		log.Info("connected to docker\n%#v", version)
	}
	d.docker.StartMonitorEvents(d.eventCallback, nil, name, handler)
}

func (d DockerMonitor) Stop() {
	log.Info("stoping docker events monitoring")
	d.docker.StopAllMonitorEvents()
}

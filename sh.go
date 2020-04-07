package main

var shellTemplate = `cd {{ .BaseDir }}
docker run -d --name ethereum-node -v $PWD:/root ethereum/client-go init {{.Genesis}}
sleep 1
docker stop ethereum-node
docker rm ethereum-node

{{range .Nodes}}sudo cp -r .ethereum  {{ . }} {{"\n"}} {{end}}
docker-compose -f {{.Compose}} up -d
`

// Shell is used for shell template to generate shell to start ethereum network.
type Shell struct {
	BaseDir string
	Genesis string
	Nodes   []string
	Compose string
}

package main

var shellTemplate = `cd {{ .BaseDir }}
docker run -d --name ethereum-node -v {{ .BaseDir }}:/root ethereum/client-go init {{.Genesis}}
sleep 1
docker stop ethereum-node
docker rm ethereum-node

{{range .Nodes}}sudo cp -r .ethereum  {{ . }} {{"\n"}} {{end}}
sudo rm -fr .ethereum
docker-compose -f {{.Compose}} up -d
`

// Shell is used for shell template to generate shell to start ethereum network.
type Shell struct {
	BaseDir string
	Genesis string
	Nodes   []string
	Compose string
}

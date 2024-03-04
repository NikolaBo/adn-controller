package controllers

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type AppManifest struct {
	AppName      string                       `yaml:"app_name"`
	AppStructure []string                     `yaml:"app_structure"`
	Edge         map[string][]EdgeElementItem `yaml:"edge"`
	Link         map[string][]PairElementItem `yaml:"link"`
}

type EdgeElementItem struct {
	Method   string `yaml:"method"`
	Name     string `yaml:"name"`
	Position string `yaml:"position"`
	Proto    string `yaml:"proto"`
}

type PairElementItem struct {
	Method   string `yaml:"method"`
	Name1    string `yaml:"name1"`
	Name2    string `yaml:"name2"`
	Position string `yaml:"position"`
	Proto    string `yaml:"proto"`
}

func removeParentheses(s string) string {
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	return s
}

func ConvertToADNSpec(appName, clientService, serverService, clientChain, serverChain, anyChain, pairChain, method, proto, fileName string) error {
	clientServerTag := fmt.Sprintf("%s->%s", clientService, serverService)

	appManifest := AppManifest{
		AppName: appName,
		AppStructure: []string{
			clientServerTag,
		},
		Edge: make(map[string][]EdgeElementItem),
		Link: make(map[string][]PairElementItem),
	}

	if len(clientChain) > 0 {
		client_elements := strings.Split(clientChain, "->")
		for _, element := range client_elements {
			appManifest.Edge[clientServerTag] = append(appManifest.Edge[clientServerTag], EdgeElementItem{
				Method:   method,
				Name:     removeParentheses(element),
				Position: "C",
				Proto:    proto,
			})
		}
	}

	if len(serverChain) > 0 {
		server_elements := strings.Split(serverChain, "->")
		for _, element := range server_elements {
			appManifest.Edge[clientServerTag] = append(appManifest.Edge[clientServerTag], EdgeElementItem{
				Method:   method,
				Name:     removeParentheses(element),
				Position: "S",
				Proto:    proto,
			})
		}
	}

	if len(anyChain) > 0 {
		any_elements := strings.Split(anyChain, "->")
		for _, element := range any_elements {
			appManifest.Edge[clientServerTag] = append(appManifest.Edge[clientServerTag], EdgeElementItem{
				Method:   method,
				Name:     removeParentheses(element),
				Position: "C/S",
				Proto:    proto,
			})
		}
	}

	// TODO: Update pairChain parsing logic
	if len(pairChain) > 0 {
		pair_elements := strings.Split(pairChain, "->")
		for _, element := range pair_elements {
			// Modify as needed to specify different elements for Name1 and Name2
			appManifest.Link[clientServerTag] = append(appManifest.Link[clientServerTag], PairElementItem{
				Method: method,
				Name1:  removeParentheses(element),
				Name2:  removeParentheses(element), // Modify as needed
				Proto:  proto,
			})
		}
	}

	yamlBytes, err := yaml.Marshal(&appManifest)
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(yamlBytes)
	if err != nil {
		return err
	}

	return nil
}

// func main() {
// 	clientService := "frontend"
// 	serverService := "ping"
// 	name := "ping"
// 	method := "PingEcho"
// 	clientChain := "rate()->fault()"
// 	serverChain := "admission()->loadbalance()"
// 	anyChain := "metrics()->tracing()"
// 	pairChain := "encrypt()->decrypt()"
// 	proto := "/users/xzhu/adn-compiler/examples/proto/ping.proto"
// 	fileName := "output.yaml"

// 	err := GenerateAndWriteYAMLToFile(clientService, serverService, clientChain, serverChain, anyChain, pairChain, name, method, proto, fileName)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("YAML written to %s\n", fileName)
// }
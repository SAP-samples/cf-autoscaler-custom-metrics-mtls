package app_config

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfenv"
)

/*
  Mocking package level functions
*/
/*
	Technique#1 - Higher order function - use to mock the package level functions
*/
type currentAppEnv func() (*cfenv.App, error)

/*
	Technique#2
	Monkey Patching - use to mock the package level functions - cfenv.IsRunningOnCF
*/
var IsRunningOnCF = cfenv.IsRunningOnCF

func GetAppEnv(currentEnv currentAppEnv) (*cfenv.App, error) {
	if !IsRunningOnCF() {
		return &cfenv.App{}, fmt.Errorf("cloud Foundry env not found.Please run the application of CF env")
	}
	appEnv, _ := currentEnv()
	fmt.Println("ID:", appEnv.ID)
	fmt.Println("Index:", appEnv.Index)
	fmt.Println("Name:", appEnv.Name)
	fmt.Println("Host:", appEnv.Host)
	fmt.Println("Port:", appEnv.Port)
	fmt.Println("Version:", appEnv.Version)
	fmt.Println("Home:", appEnv.Home)
	fmt.Println("MemoryLimit:", appEnv.MemoryLimit)
	fmt.Println("WorkingDir:", appEnv.WorkingDir)
	fmt.Println("TempDir:", appEnv.TempDir)
	fmt.Println("User:", appEnv.User)
	fmt.Println("Services:", appEnv.Services)
	return appEnv, nil

}

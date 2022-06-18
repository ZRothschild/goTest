package main

import (
	"github.com/casbin/casbin/v2"
	"log"
	"os"
)

var (
	logger *log.Logger
)

func init() {
	//O_APPEND 添加写  O_CREATE 不存在则生成   O_WRONLY 只写模式
	f, err := os.OpenFile("./text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	//defer func(f *os.File) {
	//	if err = f.Close(); err != nil {
	//		log.Fatalln(err)
	//
	//	}
	//}(f)
	//时间展示格式  LstdFlags
	logger = log.New(f, "prefix ", log.LstdFlags)
}

func main() {
	var (
		err error
		e   *casbin.Enforcer
		sub = "name" // the user that wants to access a resource.
		dom = "domain1"
		obj = "data3" // the resource that is going to be accessed.
		act = "read"  // the operation that the user performs on the resource.
	)

	if e, err = casbin.NewEnforcer("./rbac_with_domains_model.conf", "rbac_with_domains_policy.csv"); err != nil {
		logger.Printf("casbin.NewEnforcer err %v", err)
		return
	}
	policy := e.GetPolicy()
	logger.Println(policy)

	groupingPolicy := e.GetGroupingPolicy()
	logger.Println(groupingPolicy)

	allSubjects := e.GetAllSubjects()
	logger.Println(allSubjects)

	if res, _ := e.Enforce(sub, dom, obj, act); res {
		logger.Print("casbin.NewEnforcer 通过")
	} else {
		logger.Print("casbin.NewEnforcer 不通过")
	}

	roles, _ := e.GetImplicitRolesForUser("name", "domain1")
	logger.Printf("casbin.NewEnforcer roles %v", roles)

}

package main

import (
	"flag"
	"fmt"
)

func main() {
	false_dir := flag.String("false_dir", "", "enter the false data dir path")
	right_dir := flag.String("true_dir", "", "enter the right data dir path")
	flag.Parse() // 解析参数
	if len(*false_dir) == 0 || len(*right_dir) == 0 {
		panic("unexpected path")
	}
	fmt.Printf("%s:%s\n", *false_dir, *right_dir)
	moduleList, rightList, falseList := findDiffModule(*false_dir, *right_dir)

	fmt.Println(moduleList)

	for i := 0; i < len(moduleList); i++ {
		fmt.Printf("------------module %s begin--------------------\n", moduleList[i])
		keys, rightValues, falseValues := findDiffKeyInModule(*right_dir+rightList[i], *false_dir+falseList[i])
		fmt.Printf("total diff keys is %d\n", len(keys))
		for j := 0; j < len(keys); j++ {
			fmt.Printf("key is %s\n", keys[j])
			fmt.Printf("right value is %s\n", rightValues[j])
			fmt.Printf("false value is %s\n", falseValues[j])
		}
		fmt.Printf("------------module %s end  --------------------\n", moduleList[i])
	}
}

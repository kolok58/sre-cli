/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists deployments in the given namespace or defaults to all namespaces",
	Long:    ``,
	Example: "sre list --namespace [namespaceName]",
	Run: func(cmd *cobra.Command, args []string) {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// Use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// Create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		deploymentsClient := clientset.AppsV1().Deployments(namespace)
		list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		if len(list.Items) == 0 {
			fmt.Print("No deployments exist!\n")
		} else {
			for _, d := range list.Items {
				fmt.Printf(" * %s (%d replicas) Namespace: %s \n", d.Name, *d.Spec.Replicas, d.Namespace)
			}
		}
	},
}

func init() {

	listCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace to target")

	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

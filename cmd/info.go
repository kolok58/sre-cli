/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:     "info",
	Short:   "Prints info on a specified deployment",
	Long:    ``,
	Example: "sre info --deployment [deploymentName] --namespace [namespaceName]",
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

		if namespace == "" {
			// List deployments across all namespaces
			deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			namespaces := make([]string, 0)
			for _, d := range deployments.Items {
				if d.Name == deployment {
					namespaces = append(namespaces, d.Namespace)
				}
			}
			if len(namespaces) == 0 {
				fmt.Printf("Deployment %s not found\n", deployment)
				return
			} else {
				namespace = namespaces[0]
			}

		}

		deploy, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deployment, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		if deploy.Name == "" {
			fmt.Printf("Deployment %s not found\n", deployment)
			return
		}

		if full {
			deploymentJSON, err := json.MarshalIndent(deploy.Spec.Template.Spec.Containers, "", " ")
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Name: %s\nNamespace: %s\nReplicas: %d\nCreationDate: %s\nContainer Spec: %s\n",
				deploy.Name,
				deploy.Namespace,
				*deploy.Spec.Replicas,
				deploy.CreationTimestamp.Format("01-02-2006 15:04:05 MST"),
				deploymentJSON)
		} else {
			fmt.Printf("Name: %s\nNamespace: %s\nReplicas: %d\n",
				deploy.Name,
				deploy.Namespace,
				*deploy.Spec.Replicas)
		}
	},
}

func init() {
	infoCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace to target")
	infoCmd.Flags().StringVarP(&deployment, "deployment", "d", "", "The deployment to retrieve info on")
	infoCmd.Flags().BoolVar(&full, "full", false, "enable more info on the deployment")

	if err := infoCmd.MarkFlagRequired("deployment"); err != nil {
		fmt.Println(err)
	}

	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

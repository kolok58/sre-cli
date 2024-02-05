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

// scaleCmd represents the scale command
var scaleCmd = &cobra.Command{
	Use:     "scale",
	Short:   "Scales a deployment to the desired number of replicas",
	Long:    ``,
	Example: "sre scale --deployment [deploymentName] --replicas [amount] --namespace [namespaceName]",
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
			deploymentNames := 0
			for _, d := range deployments.Items {
				if d.Name == deployment {
					namespace = d.Namespace
					deploymentNames += 1
				}
			}
			if deploymentNames == 0 {
				fmt.Printf("Deployment %s not found\n", deployment)
				return
			}
			if deploymentNames >= 2 {
				fmt.Print("Namespace not specified and multiple deployments exist with the same name across various namespaces. Cannot perform scale operation\n")
				return
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

		// Update the replica count
		deploy.Spec.Replicas = &replicas

		// Update the deployment
		_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
		if err != nil {
			fmt.Printf("Failed to scale deployment: %v\n", err)
		} else {
			fmt.Printf("Scaled deployment %s to %d replicas\n", deployment, replicas)
		}
	},
}

func init() {

	scaleCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace to target")
	scaleCmd.Flags().StringVarP(&deployment, "deployment", "d", "", "The deployment to scale")
	scaleCmd.Flags().Int32VarP(&replicas, "replicas", "r", 0, "The desired number of replicas")

	if err := scaleCmd.MarkFlagRequired("deployment"); err != nil {
		fmt.Println(err)
	}
	if err := scaleCmd.MarkFlagRequired("replicas"); err != nil {
		fmt.Println(err)
	}
	rootCmd.AddCommand(scaleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scaleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scaleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

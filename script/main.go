package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type configuration struct {
	XMLName    xml.Name   `xml:"configuration"`
	Properties []property `xml:"property"`
}

type property struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

const (
	hdfsSiteXML         = "hdfs-site.xml"
	mesosSiteXML        = "mesos-site.xml"
	defaultConfDir      = "conf/default"
	xmlStylesheetHeader = "<?xml-stylesheet type=\"text/xsl\" href=\"configuration.xsl\"?>"
)

var (
	defaultHdfsConfigXML  = fmt.Sprintf("%s/%s", defaultConfDir, hdfsSiteXML)
	defaultMesosConfigXML = fmt.Sprintf("%s/%s", defaultConfDir, mesosSiteXML)
	outputDir             string
	mesosENVList          = []string{
		"MESOS_HDFS_JRE-URL",
		"MESOS_HDFS_DATA_DIR",
		"MESOS_HDFS_SECONDARY_DATA_DIR",
		"MESOS_HDFS_NATIVE-HADOOP-BINARIES",
		"MESOS_HDFS_FRAMEWORK_MNT_PATH",
		"MESOS_HDFS_STATE_ZK",
		"MESOS_MASTER_URI",
		"MESOS_HDFS_ZKFC_HA_ZOOKEEPER_QUORUM",
		"MESOS_HDFS_FRAMEWORK_NAME",
		"MESOS_HDFS_MESOSDNS",
		"MESOS_HDFS_MESOSDNS_DOMAIN",
		"MESOS_NATIVE_LIBRARY",
		"MESOS_HDFS_JOURNALNODE_COUNT",
		"MESOS_HDFS_JVM_OVERHEAD",
		"MESOS_HDFS_HADOOP_HEAP_SIZE",
		"MESOS_HDFS_NAMENODE_HEAP_SIZE",
		"MESOS_HDFS_DATANODE_HEAP_SIZE",
		"MESOS_HDFS_EXECUTOR_HEAP_SIZE",
		"MESOS_HDFS_EXECUTOR_CPUS",
		"MESOS_HDFS_NAMENODE_CPUS",
		"MESOS_HDFS_JOURNALNODE_CPUS",
		"MESOS_HDFS_DATANODE_CPUS",
		"MESOS_HDFS_DATANODE_EXCLUSIVE",
	}
)

func (conf *configuration) setProp(propName, propValue string) {
	for index, prop := range conf.Properties {
		if prop.Name == propName {
			conf.Properties[index].Value = propValue
			return
		}
	}

	conf.Properties = append(conf.Properties, property{
		Name:  propName,
		Value: propValue,
	})
}

func init() {
	flag.StringVar(&outputDir, "output_dir", "conf", "set output directory path")
}

func main() {
	flag.Parse()
	if err := createHdfsConfig(); err != nil {
		log.Fatal(err)
	}
	if err := createMesosConfig(); err != nil {
		log.Fatal(err)
	}
}

func createHdfsConfig() error {
	hdfsConf, err := readConfig(defaultHdfsConfigXML)
	if err != nil {
		return err
	}
	return writeConfig(hdfsConf, fmt.Sprintf("%s/%s", outputDir, hdfsSiteXML))
}

func createMesosConfig() error {
	mesosConf, err := readConfig(defaultMesosConfigXML)
	if err != nil {
		return err
	}
	for _, envKey := range mesosENVList {
		if envValue := os.Getenv(envKey); envValue != "" {
			propName := envKeyToPropName(envKey)
			mesosConf.setProp(propName, envValue)
		}
	}
	return writeConfig(mesosConf, fmt.Sprintf("%s/%s", outputDir, mesosSiteXML))
}

func readConfig(filePath string) (*configuration, error) {
	fBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	conf := new(configuration)
	err = xml.Unmarshal(fBytes, conf)
	return conf, err
}

func writeConfig(conf *configuration, filePath string) error {
	fBytes, err := xml.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, []byte(fmt.Sprintf("%s%s\n\n%s", xml.Header, xmlStylesheetHeader, fBytes)), 0644)
}

func envKeyToPropName(envKey string) string {
	return strings.ToLower(strings.Replace(envKey, "_", ".", -1))
}

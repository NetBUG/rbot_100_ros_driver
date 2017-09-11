package main

import (
	"encoding/json"
	"fmt"
	"github.com/CiscoZeus/zeus-analytics/registry/test/client"
)

type APIError struct {
	Message string `json:"message"`
}
type APIErrors struct {
	Errors []APIError
}

func (es *APIErrors) Append(e APIError) []APIError {
	es.Errors = append(es.Errors, e)
	return es.Errors
}

func (es *APIErrors) Len() int {
	return len(es.Errors)
}

type testBase struct {
	testClient *client.Client
}

func testAppCRUD(testClient client.Client) {
	var errors APIErrors
	createF, re := testClient.CreateApp("AnamolyDetection", "Application to detect anamolies from network logs")
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ := createF.Body.ToString()
	fmt.Printf("Create app resp: %s \n", str)

	var dat map[string]interface{}
	json.Unmarshal([]byte(str), &dat)
	appid := string(dat["id"].(string))

	readF, re := testClient.GetApp(appid)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = readF.Body.ToString()
	fmt.Printf("Read app resp: %s \n", str)

	updateF, re := testClient.UpdateApp(appid, "New Anamoly Detection", "New application to detect anomolies from network logs")
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = updateF.Body.ToString()
	fmt.Printf("Update app resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	updateAppId := string(dat["id"].(string))

	getFL, re := testClient.GetAppList()
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = getFL.Body.ToString()
	fmt.Printf("Read app list resp: %s \n", str)

	err := testClient.DeleteApp(updateAppId)
	if err != nil {
		errors.Append(APIError{Message: "Could not delete registry service."})
		fmt.Printf("Error")
		return
	}

	fmt.Printf("Deleted app: \n")

	newFL, re := testClient.GetAppList()
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = newFL.Body.ToString()
	fmt.Printf("Read app list resp: %s \n", str)
}

func testJobCRUD(testClient client.Client) {
	var errors APIErrors
	createdF, re := testClient.CreateApp("AnamolyDetection2", "Application2 to detect anamolies from network logs")
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ := createdF.Body.ToString()
	fmt.Printf("Create app resp: %s \n", str)

	var dat map[string]interface{}
	json.Unmarshal([]byte(str), &dat)
	appid := string(dat["id"].(string))

	body := client.Job{Name: "Batch training", Description: "Training algorithm", Input: "Streaming Data", Output: "Trained model", Action: "python batch.py", Parameters: "training=1000,test=20000"}
	createdJ, re := testClient.CreateJob(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create job resp: %s \n", str)

	json.Unmarshal([]byte(str), &dat)
	jobId := string(dat["id"].(string))
	readJ, re := testClient.GetJob(jobId)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = readJ.Body.ToString()
	fmt.Printf("Read job resp: %s \n", str)

	body = client.Job{Name: "Batch training 2", Description: "Training algorithm", Input: "Streaming Data", Output: "Trained model"}
	updatedJ, re := testClient.UpdateJob(jobId, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = updatedJ.Body.ToString()
	fmt.Printf("Update job resp: %s \n", str)

	newFL, re := testClient.GetAppList()
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = newFL.Body.ToString()
	fmt.Printf("Read app list resp: %s \n", str)

	/*
		deletedJ, re := testClient.DeleteJob(jobId);
		if re != nil {
			errors.Append(APIError{Message: "Could not access registry service."})
			fmt.Printf("Error")
			return
		}
		str, _ = deletedJ.Body.ToString();
		fmt.Printf("Delete job resp: %s \n", str);
	*/
	jobList, re := testClient.GetJobList(appid)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = jobList.Body.ToString()
	fmt.Printf("Get job list resp: %s \n", str)

}

func testDataCRUD(testClient client.Client) {
	var errors APIErrors
	createdF, re := testClient.CreateApp("AnamolyDetection2", "Application2 to detect anamolies from network logs")
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ := createdF.Body.ToString()
	fmt.Printf("Create app resp: %s \n", str)

	var dat map[string]interface{}
	json.Unmarshal([]byte(str), &dat)
	appid := string(dat["id"].(string))

	body := client.Data{Name: "StreamingData", Description: "Streaming data", Type: "Kafka", Parameters: "LogTopic", Store: "Kafka"}
	createdJ, re := testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)

	json.Unmarshal([]byte(str), &dat)
	dataId := string(dat["id"].(string))
	readJ, re := testClient.GetData(dataId)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = readJ.Body.ToString()
	fmt.Printf("Read data resp: %s \n", str)

	body = client.Data{Name: "StreamingData2", Description: "Streaming data", Type: "Kafka", Parameters: "LogTopic", Store: "Kafka"}
	updatedJ, re := testClient.UpdateData(dataId, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = updatedJ.Body.ToString()
	fmt.Printf("Update data resp: %s \n", str)

	deletedJ, re := testClient.DeleteData(dataId)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = deletedJ.Body.ToString()
	fmt.Printf("Delete data resp: %s \n", str)

	dataList, re := testClient.GetDataList(appid)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = dataList.Body.ToString()
	fmt.Printf("Get data list resp: %s \n", str)

}
func testScheduleCRUD(testClient client.Client) {
	var errors APIErrors
	createdF, re := testClient.CreateApp("AnamolyDetection2", "Application2 to detect anamolies from network logs")
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ := createdF.Body.ToString()
	fmt.Printf("Create app resp: %s \n", str)

	var dat map[string]interface{}
	json.Unmarshal([]byte(str), &dat)
	appid := string(dat["id"].(string))

	body := client.Schedule{Name: "Schedule", Description: "Schedule for Anamloy Detection App", Sched: "Task1, Task2"}
	createdJ, re := testClient.CreateSchedule(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create schedule resp: %s \n", str)

	json.Unmarshal([]byte(str), &dat)
	scheduleId := string(dat["id"].(string))
	readJ, re := testClient.GetSchedule(scheduleId)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = readJ.Body.ToString()
	fmt.Printf("Read schedule resp: %s \n", str)

	body = client.Schedule{Name: "Schedule", Description: "Schedule for Anamloy Detection App", Sched: "Task1, Task3"}
	updatedJ, re := testClient.UpdateSchedule(scheduleId, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = updatedJ.Body.ToString()
	fmt.Printf("Update schedule resp: %s \n", str)

	deletedJ, re := testClient.DeleteSchedule(scheduleId)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = deletedJ.Body.ToString()
	fmt.Printf("Delete schedule resp: %s \n", str)

	scheduleList, re := testClient.GetScheduleList(appid)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = scheduleList.Body.ToString()
	fmt.Printf("Get schedule list resp: %s \n", str)

}

func testFlow(testClient client.Client) {
	//Create app
	//Create IP data
	//Create OP data
	//Create Job
	//Create schedule

	//Create app
	var errors APIErrors
	createdF, re := testClient.CreateApp("NetStats", "Application to count packets statistics")
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ := createdF.Body.ToString()
	fmt.Printf("Create app resp: %s \n", str)

	var dat map[string]interface{}
	json.Unmarshal([]byte(str), &dat)
	appid := string(dat["id"].(string))

	body := client.Data{Name: "Network packets input stream", Description: "Network packet streaming data", Type: "Spark/DStream", Parameters: "", Store: "Kafka"}
	createdJ, re := testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	streamingDataId := string(dat["id"].(string))

	body = client.Data{Name: "Network stats stream", Description: "Network stats streaming data", Type: "Spark/DStream", Parameters: "", Store: "Kafka"}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	streamingAnalyticsDataId := string(dat["id"].(string))

	body = client.Data{Name: "TrainingData", Description: "Training data", Type: "HDFS", Parameters: "/user/hdpuser/sdata", Store: "HDFS"}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	trainingDataId := string(dat["id"].(string))

	body = client.Data{Name: "TrainedModel", Description: "Model built from training data", Type: "HDFS", Parameters: "/user/hdpuser/tdata", Store: "HDFS"}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	modelDataId := string(dat["id"].(string))

	//Create job
	jobbody := client.Job{Name: "Training", Description: "Training algorithm", Input: trainingDataId, Output: modelDataId, Action: "python training.py", Parameters: "training=1000,test=20000"}
	createdJ, re = testClient.CreateJob(appid, jobbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create job resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	trainingJobId := string(dat["id"].(string))

	jobbody = client.Job{Name: "Analytics", Description: "Streaming data analytics algorithm", Input: modelDataId + "," + streamingDataId, Output: streamingAnalyticsDataId, Action: "python streaming.py", Parameters: "intval=1000"}
	createdJ, re = testClient.CreateJob(appid, jobbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create job resp: %s \n", str)

	json.Unmarshal([]byte(str), &dat)
	streamingAnalyticsId := string(dat["id"].(string))

	//Create schedule
	schedbody := client.Schedule{Name: "Schedule", Description: "Schedule for Anamoly Detection App", Sched: trainingJobId + "," + streamingAnalyticsId}
	createdJ, re = testClient.CreateSchedule(appid, schedbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create schedule resp: %s \n", str)

	//Check everything
	getFL, re := testClient.GetAppList()
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = getFL.Body.ToString()
	fmt.Printf("Get job list resp: %s \n", str)

	jobList, re := testClient.GetJobList(appid)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = jobList.Body.ToString()
	fmt.Printf("Get job list resp: %s \n", str)

	dataList, re := testClient.GetDataList(appid)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = dataList.Body.ToString()
	fmt.Printf("Get data list resp: %s \n", str)

	scheduleList, re := testClient.GetScheduleList(appid)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = scheduleList.Body.ToString()
	fmt.Printf("Get schedule list resp: %s \n", str)

	//Clone app
	cloneF, re := testClient.CloneApp(appid)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = cloneF.Body.ToString()
	fmt.Printf("Clone app resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)

	getFL, re = testClient.GetAppList()
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = getFL.Body.ToString()

}

type HParameters struct {
	Uri      string `json:"uri"`
	Username string `json:"username"`
}
type Parameters struct {
	Source     string `json:"source"`
	Topic     string `json:"topic"`
}
type Jobparameters struct {
	Mainclass string `json:"mainclass"`
	Jars string `json:"jars"`
	Resource       string `json:"resource"`
	Args      string `json:"args"`
	Secs      string `json:"secs"`
}
type Jobparameters2 struct {
	Cmd string `json:"cmd"`
	Image       string `json:"image"`
	Args      string `json:"args"`
}

func testRealApp0(testClient client.Client) {
	//Create app
	//Create IP data
	//Create OP data
	//Create Job
	//Create schedule

	//Create app
	var errors APIErrors
	createdF, re := testClient.CreateApp("NetAnalytics", "Network traffic statistics and predictive analytics")

	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ := createdF.Body.ToString()
	fmt.Printf("Create app resp: %s \n", str)

	var dat map[string]interface{}
	json.Unmarshal([]byte(str), &dat)
	appid := string(dat["id"].(string))

	params := Parameters{Source: "netpkts", Topic: ""}
	ps, _ := json.Marshal(&params)
	body := client.Data{Name: "Network packets", Description: "Network traffic stream", Type: "Spark/DStream",  Parameters: string(ps)}
	createdJ, re := testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	streamingDataId := string(dat["id"].(string))

	params = Parameters{Source: "", Topic: "analysis"}
	ps, _ = json.Marshal(&params)
	body = client.Data{Name: "Network analysis", Description: "Network stats/analysis stream", Type: "Spark/DStream", Parameters: string(ps)}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	streamingAnalyticsDataId := string(dat["id"].(string))

	hparams := HParameters{Username: "hdpuser", Uri: "hdfs://localhost:9000"}
	ps, _ = json.Marshal(&hparams)
	body = client.Data{Name: "Training data", Description: "Training data for network traffic prediction", Type: "HDFS", Parameters: string(ps)}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	trainingDataId := string(dat["id"].(string))

	hparams = HParameters{Username: "hdpuser", Uri: "hdfs://localhost:9000"}
	ps, _ = json.Marshal(&hparams)
	body = client.Data{Name: "Model", Description: "Trained model for network traffic prediction", Type: "HDFS", Parameters: string(ps)}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	modelDataId := string(dat["id"].(string))

	//Create streaming job
	jparams := Jobparameters{Mainclass: "com.ciscozeusdev.apps.networkanalytics.NetAnalytics", Resource: "/home/ubuntu/zeus-analytics-apps/networkanalytics/target/scala-2.11/ciscozeusdev-apps-networkanalytics-assembly-0.1.jar", Args: "",Secs: "10" }
	jps, _ := json.Marshal(&jparams)
	jobbody := client.Job{Name: "NetAnalytics", Description: "Statistics and predictive analytics on streaming network traffic", Type: "Spark/Streaming", Input: streamingDataId + "," + modelDataId, Output: streamingAnalyticsDataId + "," + trainingDataId, Action: "Continuous", Parameters: string(jps)}
	createdJ, re = testClient.CreateJob(appid, jobbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create job resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	statsJobId := string(dat["id"].(string))

	jparams = Jobparameters{Mainclass: "com.ciscozeusdev.apps.networkanalytics.NetAnalyticsTrain", Resource: "/home/ubuntu/zeus-analytics-apps/networkanalytics/target/scala-2.11/ciscozeusdev-apps-networkanalytics-assembly-0.1.jar", Args: "", Secs: "10"}
	jps, _ = json.Marshal(&jparams)
	jobbody = client.Job{Name: "NetAnalyticsTrain", Description: "Training for prediction on streaming network traffic", Type: "Spark/Batch", Input: trainingDataId, Output: modelDataId, Action: "One Time", Parameters: string(jps)}
	createdJ, re = testClient.CreateJob(appid, jobbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create job resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	statsJobId = string(dat["id"].(string))
	
        //Create streaming job
	jparams2 := Jobparameters2{Cmd: "",  Image: "networkanalytics/visualize", Args: ""}
	jps, _ = json.Marshal(&jparams2)
	jobbody = client.Job{Name: "Visualize", Description: "Network  analytics visualization", Type: "Container", Input: streamingAnalyticsDataId , Output: "" , Action: "Continuous", Parameters: string(jps)}
	createdJ, re = testClient.CreateJob(appid, jobbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create job resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	statsJobId = string(dat["id"].(string))


	//Create schedule
	schedbody := client.Schedule{Name: "Schedule", Description: "Schedule for Network Traffic Analytics App", Sched: statsJobId}
	createdJ, re = testClient.CreateSchedule(appid, schedbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create schedule resp: %s \n", str)

}
func testRealApp1(testClient client.Client) {
	//Create app
	//Create IP data
	//Create OP data
	//Create Job
	//Create schedule

	//Create app
	var errors APIErrors
	createdF, re := testClient.CreateApp("SentimentAnalysis", "Sentiment and topic analysis of conversations")

	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ := createdF.Body.ToString()
	fmt.Printf("Create app resp: %s \n", str)

	var dat map[string]interface{}
	json.Unmarshal([]byte(str), &dat)
	appid := string(dat["id"].(string))

	params := Parameters{Source: "", Topic: "conversationstream"}
	ps, _ := json.Marshal(&params)
	body := client.Data{Name: "Conversation stream", Description: "Conversation stream", Type: "Spark/DStream",  Parameters: string(ps)}
	createdJ, re := testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	conversationStreamDataId := string(dat["id"].(string))
	
	params = Parameters{}
	ps, _ = json.Marshal(&params)
	body = client.Data{Name: "REST input", Description: "REST API", Type: "Port", Parameters: string(ps)}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	restDataId := string(dat["id"].(string))

	params = Parameters{Source: "", Topic: "conversationanalysis"}
	ps, _ = json.Marshal(&params)
	body = client.Data{Name: "Conversation analysis", Description: "Conversation analysis", Type: "Spark/DStream", Parameters: string(ps)}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	conversationAnalysisDataId := string(dat["id"].(string))


	params = Parameters{Source: "", Topic: "conversationstream"}
	ps, _ = json.Marshal(&params)
	body = client.Data{Name: "Conversation stream", Description: "Conversation stream", Type: "Kafka",  Parameters: string(ps)}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	conversationStreamKDataId := string(dat["id"].(string))

	params = Parameters{Source: "", Topic: "conversationanalysis"}
	ps, _ = json.Marshal(&params)
	body = client.Data{Name: "Conversation analysis", Description: "Conversation analysis", Type: "Kafka", Parameters: string(ps)}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	conversationAnalysisKDataId := string(dat["id"].(string))

	//Create streaming job
	jparams := Jobparameters{Jars: "/home/ubuntu/zeus-analytics-apps/sentimentanalysis/src/stanford-corenlp-3.6.0-models.jar",  Mainclass: "com.ciscozeusdev.apps.sentimentanalysis.SentimentAnalysis", Resource: "/home/ubuntu/zeus-analytics-apps/sentimentanalysis/target/scala-2.11/ciscozeusdev-apps-sentimentanalysis-assembly-0.1.jar", Args: "", Secs: "10"}
	jps, _ := json.Marshal(&jparams)
	jobbody := client.Job{Name: "SentimentAnalysis", Description: "Sentiment Analysis", Type: "Spark/Streaming", Input: conversationStreamDataId, Output: conversationAnalysisDataId, Action: "Continuous", Parameters: string(jps)}
	createdJ, re = testClient.CreateJob(appid, jobbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create job resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	statsJobId := string(dat["id"].(string))

	jparams = Jobparameters{Mainclass: "com.ciscozeusdev.apps.sentimentanalysis.AlexaServer", Resource: "/home/ubuntu/zeus-analytics-apps/sentimentanalysis/target/scala-2.11/ciscozeusdev-apps-sentimentanalysis-assembly-0.1.jar", Args: "", Secs: "10"}
	jps, _ = json.Marshal(&jparams)
	jobbody = client.Job{Name: "AlexaServer", Description: "Alexa Server", Type: "Spark/Streaming", Input: conversationAnalysisKDataId+","+restDataId, Output: conversationStreamKDataId, Action: "Continuous", Parameters: string(jps)}
	createdJ, re = testClient.CreateJob(appid, jobbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create job resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	statsJobId = string(dat["id"].(string))

	//Create schedule
	schedbody := client.Schedule{Name: "Schedule", Description: "Schedule for Network Traffic Analytics App", Sched: statsJobId}
	createdJ, re = testClient.CreateSchedule(appid, schedbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create schedule resp: %s \n", str)

}



func testRealApp2(testClient client.Client) {
	//Create app
	//Create IP data
	//Create OP data
	//Create Job
	//Create schedule

	//Create app
	var errors APIErrors
	createdF, re := testClient.CreateApp("VideoEmotionDetection", "Video emotion detection")

	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ := createdF.Body.ToString()
	fmt.Printf("Create app resp: %s \n", str)

	var dat map[string]interface{}
	json.Unmarshal([]byte(str), &dat)
	appid := string(dat["id"].(string))

	params := Parameters{Source: "", Topic: "frames"}
	ps, _ := json.Marshal(&params)
	body := client.Data{Name: "Frames", Description: "Video frames", Type: "Kafka",  Parameters: string(ps)}
	createdJ, re := testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	framesDataId := string(dat["id"].(string))

	params = Parameters{Source: "", Topic: "emotions"}
	ps, _ = json.Marshal(&params)
	body = client.Data{Name: "Emotions", Description: "Emotions detected from video", Type: "Kafka", Parameters: string(ps)}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}

	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	emotionsDataId := string(dat["id"].(string))
	
	params = Parameters{}
	ps, _ = json.Marshal(&params)
	body = client.Data{Name: "Web server input", Description: "Web server", Type: "Port", Parameters: string(ps)}
	createdJ, re = testClient.CreateData(appid, body)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create data resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	webserverDataId := string(dat["id"].(string))
	

	//Create streaming job
	jparams := Jobparameters2{Cmd: "",  Image: "ciscozeus/emotiondetection", Args: ""}
	jps, _ := json.Marshal(&jparams)
	jobbody := client.Job{Name: "EmotionDetection", Description: "Emotion detection from frame stream", Type: "Container", Input: framesDataId, Output: emotionsDataId, Action: "Continuous", Parameters: string(jps)}
	createdJ, re = testClient.CreateJob(appid, jobbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create job resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	statsJobId := string(dat["id"].(string))

	jparams = Jobparameters2{Cmd: "",  Image: "ciscozeus/framedetection", Args: ""}
	jps, _ = json.Marshal(&jparams)
	jobbody = client.Job{Name: "FrameDetection", Description: "Frame detection from video stream", Type: "Container", Input: emotionsDataId+","+webserverDataId, Output: framesDataId, Action: "Continuous", Parameters: string(jps)}
	createdJ, re = testClient.CreateJob(appid, jobbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create job resp: %s \n", str)
	json.Unmarshal([]byte(str), &dat)
	statsJobId = string(dat["id"].(string))
	
	//Create schedule
	schedbody := client.Schedule{Name: "Schedule", Description: "Schedule for Network Traffic Analytics App", Sched: statsJobId}
	createdJ, re = testClient.CreateSchedule(appid, schedbody)
	if re != nil {
		errors.Append(APIError{Message: "Could not access registry service."})
		fmt.Printf("Error")
		return
	}
	str, _ = createdJ.Body.ToString()
	fmt.Printf("Create schedule resp: %s \n", str)

}

func main() {
	testClient := client.Client{Host: "localhost", Port: "8800"}
	//testAppCRUD(testClient);
	//testJobCRUD(testClient);
	//testDataCRUD(testClient);
	//testScheduleCRUD(testClient);
	//testFlow(testClient);
	testRealApp0(testClient)
	testRealApp1(testClient)
	testRealApp2(testClient)

}

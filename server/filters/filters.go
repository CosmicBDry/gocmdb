package filters

import (
	"time"

	"strconv"

	"github.com/astaxie/beego/context"
	"github.com/prometheus/client_golang/prometheus"
)

//定义标签------------------------------------------------------------------------>
var (
	//Counter计数指标类型
	totalRequest = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cmdb_request_total",
		Help: "",
	})
	//CounterVec带可变标签的计数指标类型
	urlRequest = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cmdb_url_request_total",
		Help: "",
	}, []string{"url"})
	//CounterVec带可变标签的计数指标类型
	statusCode = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cmdb_status_code",
		Help: "",
	}, []string{"code"})
	//Histogram直方图指标类型，统计各个桶区间出现样本数据的次数，次数即为样本值，为累计统计即某个桶区间包含前面所有桶区间的样本值
	elapsedTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{

		Name: "cmdb_request_url_elapsetime",
		Help: "",
	}, []string{"url"})
)

//标签值更新------------------------------------------------------------------------>
//之所以定义为此类型的函数BeforeExec(*context.Context)，是提供于beego.InsertFilter()过滤器中作为处理的方法函数
//BeforeExec代表执行之前的过滤器操作，此过滤器就是BeforeExec
func BeforeExec(ctx *context.Context) {
	totalRequest.Inc()
	urlRequest.WithLabelValues(ctx.Input.URL()).Inc()
	//注意SetData设置的key-value保存在ctx中，只能从ctx中获取
	ctx.Input.SetData("startTime", time.Now())
}

//AfterExec代表执行之后的过滤器操作
func AfterExec(ctx *context.Context) {
	statusCode.WithLabelValues(strconv.Itoa(ctx.ResponseWriter.Status)).Inc()
	//ctx.Input.GetData获取ctx中设置的key-value
	if stime := ctx.Input.GetData("startTime"); stime != nil {
		s, ok := stime.(time.Time)
		if ok {
			//通过Sub()减去初始时间获取时间差值，即一次请求的响应时间
			elapstime := time.Now().Sub(s)
			//elapstime.Seconds()时间单位转换为秒，转化后的数据类型为float64
			elapsedTime.WithLabelValues(ctx.Input.URL()).Observe(elapstime.Seconds())
		}
	}
}

func init() {
	prometheus.MustRegister(totalRequest, urlRequest, statusCode, elapsedTime)

}

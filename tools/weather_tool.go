package tools

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type GetWeatherInput struct {
	City string `json:"city_name"`
}

func GetWeather(_ context.Context, input *GetWeatherInput) (string, error) {
	return fmt.Sprintf("%s 的天气晴，温度25摄氏度", input.City), nil
}

func CreateWeatherTool() tool.InvokableTool {
	getWeatherTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_weather",
		Desc: "get weather by city name",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"name": &schema.ParameterInfo{
					Type:     schema.String,
					Desc:     "city's name",
					Required: true,
				},
			},
		),
	}, GetWeather)
	return getWeatherTool
}

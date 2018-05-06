package server

import (
// "encoding/json"
// "sort"
)

type DownloadFlow struct {
	Time  string    `json:"time"`
	Value FlowValue `json:"values"`
}

type FlowValue struct {
	Flow int64 `json:"flow"`
}

type FlowWrapper struct {
	flows []DownloadFlow
	by    func(p, q *DownloadFlow) bool
}

func (fw FlowWrapper) Len() int {
	return len(fw.flows)
}

func (fw FlowWrapper) Swap(i, j int) {
	fw.flows[i], fw.flows[j] = fw.flows[j], fw.flows[i]
}

func (fw FlowWrapper) Less(i, j int) bool { // 重写 Less() 方法
	return fw.by(&fw.flows[i], &fw.flows[j])
}

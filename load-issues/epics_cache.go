package main

import (
	"github.com/fredex42/mm-jira-migration/common"
	"log"
)

type EpicsCache struct {
	KnownEpics map[string]string
}

func NewEpicsCache(hostname *string, jira *common.ScriptKey, pageSize int) (*EpicsCache, error) {
	epicsList, err := common.SyncLoadAllEpics(*hostname, jira, pageSize)
	if err != nil {
		log.Fatal("ERROR Could not load in epics: ", err)
	}

	out := EpicsCache{KnownEpics: make(map[string]string, len(epicsList))}

	log.Printf("INFO Loaded %d epics", len(epicsList))

	for _, e := range epicsList {
		if e.Fields.EpicName != nil {
			out.KnownEpics[e.Key] = *e.Fields.EpicName
		}
	}
	return &out, nil
}

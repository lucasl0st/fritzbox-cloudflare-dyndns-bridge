package main

import (
	"context"
	"errors"
	"github.com/cloudflare/cloudflare-go"
)

func updateDomainRecord(ctx context.Context, user User, zoneId string, address string, recordName string, zoneName string, t string) error {
	recordId, err := getRecordId(user, zoneId, recordName, zoneName, t)

	if err != nil {
		return err
	}

	return user.cloudflare.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneId), cloudflare.UpdateDNSRecordParams{
		Type:    t,
		Content: address,
		TTL:     1,
		ID:      recordId,
	})
}

func getZoneId(user User, name string) (string, error) {
	zones, err := user.cloudflare.ListZones(context.Background())

	if err != nil {
		return "", err
	}

	for _, zone := range zones {
		if zone.Name == name {
			return zone.ID, nil
		}
	}

	return "", errors.New("could not find zone")
}

func getRecordId(user User, zoneId string, name string, zoneName string, t string) (string, error) {
	records, _, err := user.cloudflare.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zoneId), cloudflare.ListDNSRecordsParams{})

	if err != nil {
		return "", err
	}

	for _, record := range records {
		if record.Name == name+zoneName && record.Type == t {
			return record.ID, nil
		}
	}

	return "", errors.New("could not find record")
}

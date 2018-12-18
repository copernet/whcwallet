package model

import (
	"database/sql"
	"errors"
	"math"
	"strconv"

	"github.com/copernet/whccommon/model"
	"github.com/copernet/whcwallet/model/cache"
	"github.com/copernet/whcwallet/util"
	"github.com/jinzhu/gorm"
)

type PropertyWithTime struct {
	PropertyData string `json:"property_data"`
	BlockTime    int64  `json:"block_time"`
}

func ListAllProperties(keyword string, pageSize, pageNo int) ([]PropertyWithTime, error) {
	var rows *sql.Rows
	var err error

	if keyword != "" {
		pid, err := strconv.Atoi(keyword)
		if err != nil {
			rows, err = db.Table("smart_properties").
				Select("smart_properties.property_data, txes.block_time").
				Where("property_name LIKE ?", "%"+keyword+"%").
				Joins("INNER JOIN txes ON txes.tx_id = smart_properties.create_tx_id").
				Limit(pageSize).
				Offset(pageSize*pageNo - pageSize).
				Rows()
		} else {
			rows, err = db.Table("smart_properties").
				Select("smart_properties.property_data, txes.block_time").
				Where("property_name LIKE ? OR property_id = ?", "%"+keyword+"%", pid).
				Joins("INNER JOIN txes ON txes.tx_id = smart_properties.create_tx_id").
				Limit(pageSize).
				Offset(pageSize*pageNo - pageSize).
				Rows()
		}
	} else {
		rows, err = db.Table("smart_properties").
			Select("smart_properties.property_data, txes.block_time").
			Joins("INNER JOIN txes ON txes.tx_id = smart_properties.create_tx_id").
			Limit(pageSize).
			Offset(pageSize*pageNo - pageSize).
			Rows()
	}

	if err != nil {
		return nil, err
	}

	properties := make([]PropertyWithTime, 0)
	var property PropertyWithTime
	for rows.Next() {
		err = db.ScanRows(rows, &property)
		if err != nil {
			return nil, err
		}

		properties = append(properties, property)
	}

	return properties, nil
}

func PropertyListCount(keyword string) (int, error) {
	var total int
	var err error

	if keyword != "" {
		pid, err := strconv.Atoi(keyword)
		if err != nil {
			err = db.Table("smart_properties").Where("property_name LIKE ?", "%"+keyword+"%").Count(&total).Error
		} else {
			err = db.Table("smart_properties").Where("property_name LIKE ? OR property_id = ?", "%"+keyword+"%", pid).Count(&total).Error
		}
	} else {
		err = db.Table("smart_properties").Count(&total).Error
	}

	if err != nil {
		return 0, err
	}

	return total, nil
}

type PropertiesByAddresses struct {
	TxData       string `json:"tx_data"`
	PropertyData string `json:"property_data"`
}

func GetPropertiesByAddresses(addresses []string, pageSize, pageNo int) ([]PropertiesByAddresses, error) {
	rows, err := db.Raw("SELECT j.tx_data, sp.property_data from tx_jsons as j INNER JOIN smart_properties as sp "+
		"ON j.tx_id = sp.create_tx_id WHERE sp.protocol <> ? AND sp.issuer IN (?) ORDER BY property_id desc",
		model.Fiat, addresses).
		Limit(pageSize).
		Offset(pageSize*pageNo - pageSize).
		Rows()

	if err != nil {
		return nil, err
	}

	sps := make([]PropertiesByAddresses, 0)
	var txdata, propertydata string
	for rows.Next() {
		err = rows.Scan(&txdata, &propertydata)
		if err != nil {
			return nil, err
		}
		sps = append(sps, PropertiesByAddresses{
			TxData:       txdata,
			PropertyData: propertydata,
		})
	}

	return sps, nil
}

func GetPropertiesNumberByAddresses(addresses []string) (int, error) {
	var count int
	row := db.Raw("SELECT count(*) as count from tx_jsons as j INNER JOIN smart_properties as sp "+
		"ON j.tx_id = sp.create_tx_id WHERE sp.protocol <> ? AND sp.issuer IN (?) ORDER BY property_id",
		model.Fiat, addresses).Row()
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetPropertyById(pid uint64) (string, error) {
	var property model.SmartProperty
	notFound := db.Table("smart_properties").Where("property_id = ?", pid).
		Select("property_data").First(&property).RecordNotFound()
	if !notFound {
		return property.PropertyData, nil
	}

	return "", gorm.ErrRecordNotFound
}

func GetPropertyByName(name string) ([]string, error) {
	var properties []*model.SmartProperty
	notFound := db.Table("smart_properties").Where("property_name = ?", name).
		Select("property_data").Find(&properties).RecordNotFound()
	if !notFound {
		ret := make([]string, 0, len(properties))
		for i := 0; i < len(properties); i++ {
			ret = append(ret, properties[i].PropertyData)
		}
		return ret, nil
	}

	return nil, gorm.ErrRecordNotFound
}

// support query via property name or property id input
func GetPropertyByKeyword(keyword string) ([]string, error) {
	var properties []*model.SmartProperty
	notFound := db.Table("smart_properties").Where("property_name LIKE ?", keyword+"%").
		Or("property_id = ?", keyword).
		Select("property_data").
		Limit(10).
		Find(&properties).
		RecordNotFound()

	if !notFound {
		ret := make([]string, 0, len(properties))
		for i := 0; i < len(properties); i++ {
			ret = append(ret, properties[i].PropertyData)
		}
		return ret, nil
	}

	return nil, gorm.ErrRecordNotFound
}

type AddressPropertyListWithBalance struct {
	Balance  float64     `json:"balance"`
	Property interface{} `json:"property"`
}

func GetPropertyListByAddress(addr string) ([]AddressPropertyListWithBalance, error) {
	rows, err := db.Raw("select ab.balance_available, sp.property_data from address_balances as ab "+
		"INNER JOIN smart_properties as sp "+
		"ON ab.property_id = sp.property_id "+
		"WHERE ab.address = ?", addr).Rows()
	if err != nil {
		return nil, err
	}

	var ret []AddressPropertyListWithBalance
	var balance float64
	var property string
	for rows.Next() {
		err := rows.Scan(&balance, &property)
		if err != nil {
			return nil, err
		}

		sp := util.JsonStringToMap(property)
		precision, ok := sp["precision"]
		if !ok {
			return nil, errors.New("property data not contain 'precision' field")
		}

		var precisionInt int
		switch p := precision.(type) {
		case string:
			precisionInt, err = strconv.Atoi(p)
			if err != nil {
				return nil, errors.New("convert precision from string to integer filed")
			}
		case float64:
			precisionInt = int(p)
		}

		balance = balance / math.Pow10(precisionInt)
		ret = append(ret, AddressPropertyListWithBalance{
			Balance:  balance,
			Property: util.JsonStringToMap(property),
		})
	}

	return ret, nil
}

func GetPropertyName(pid uint64) (string, error) {
	// The pending created issuance transaction, the property name is not set;
	// The invalid created issuance transaction, the property id = 0;
	// Bitcoin cash transaction in addresses_in_txes table, the property id = 0;
	if pid == 0 {
		return "", nil
	}

	// get property name from cache firstly.
	name, ok := cache.GetPNameCache().GetPropertyName(int64(pid))
	if ok {
		return name, nil
	}

	var sp model.SmartProperty
	err := db.Select("property_name").Where("property_id = ?", pid).First(&sp).Error
	if err != nil {
		return "", err
	}

	return sp.PropertyName, nil
}

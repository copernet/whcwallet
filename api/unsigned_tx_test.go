package api_test

import (
	"net/url"
	"testing"
)

func checkRequest(fun func() (string, url.Values), t *testing.T) {
	uri, data := fun()
	t.Log(uri)
	res := checkPost(uri, data, t)
	result, ok := res.Result.([]interface{})
	if !ok {
		t.Log("result format error")
	}
	t.Log(result)
}

func TestGetUnsignedTx(t *testing.T) {
	checkRequest(createManagedIssuanceTransaction, t)
	checkRequest(freezeTokenTransaction, t)

	checkRequest(createSimpleSendTransaction, t)
	checkRequest(createParticipateCrowdSaleTransaction, t)
	checkRequest(createSendStoTransaction, t)
	checkRequest(createSendAllTransaction, t)

	checkRequest(createFixedIssuanceTransaction, t)
	checkRequest(createCrowdSaleIssuanceTransaction, t)
	checkRequest(createCloseCrowdSaleTransaction, t)
	checkRequest(createSendGrantTransaction, t)
	checkRequest(createGetWhcTransaction, t)
	checkRequest(createChangeIssuerTransaction, t)


}

func createSimpleSendTransaction() (string, url.Values) {
	uri := "/getunsigned/0"

	result := url.Values{
		"transaction_version": {"0"},
		"pubkey":              {""},
		"fee":                 {"0.001"},
		"transaction_from":    {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"transaction_to":      {"bchtest:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk"},
		"currency_identifier": {"1"},
		"amount_to_transfer":  {"1.001"},
		"redeem_address":      {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"reference_amount":    {""},
	}

	return uri, result
}

func createParticipateCrowdSaleTransaction() (string, url.Values) {
	uri := "/getunsigned/1"
	result := url.Values{
		"transaction_version": {"0"},
		"pubkey":              {""},
		"fee":                 {"0.001"},
		"transaction_from":    {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"transaction_to":      {"bchtest:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk"},
		"currency_identifier": {"1"},
		"amount_to_transfer":  {"1.001"},
		"redeem_address":      {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"reference_amount":    {""},
	}

	return uri, result
}

func createSendStoTransaction() (string, url.Values) {
	uri := "/getunsigned/3"

	result := url.Values{
		"transaction_version":  {"0"},
		"pubkey":               {""},
		"fee":                  {"0.001"},
		"transaction_from":     {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"propertyid":           {"1"},
		"amount":               {"1"},
		"redeem_address":       {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"distributionproperty": {"1"},
	}

	return uri, result
}

func createSendAllTransaction() (string, url.Values) {
	uri := "/getunsigned/4"

	result := url.Values{
		"transaction_version": {"0"},
		"pubkey":              {""},
		"fee":                 {"0.001"},
		"transaction_from":    {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"transaction_to":      {"bchtest:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk"},
		"ecosystem":           {"1"},
		"redeem_address":      {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"reference_amount":    {""},
	}

	return uri, result
}

func createFixedIssuanceTransaction() (string, url.Values) {
	uri := "/getunsigned/50"

	result := url.Values{
		"transaction_version": {"0"},
		"pubkey":              {""},
		"fee":                 {"0.001"},
		"transaction_from":    {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"ecosystem":           {"1"},
		"precision":           {"1"},
		//"previous_property_id": {"1"},
		"property_category":    {"shopping"},
		"property_subcategory": {"supermarket"},
		"property_name":        {"test"},
		"property_url":         {"http://test.com"},
		"property_data":        {"test value"},
		"number_properties":    {"2.1"},
	}

	return uri, result
}

func createCrowdSaleIssuanceTransaction() (string, url.Values) {
	uri := "/getunsigned/51"

	result := url.Values{
		"transaction_version": {"0"},
		"pubkey":              {""},
		"fee":                 {"0.001"},
		"transaction_from":    {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"ecosystem":           {"1"},
		"precision":           {"1"},
		//"previous_property_id": {""},
		"property_category":           {"shopping"},
		"property_subcategory":        {"supermarket"},
		"property_name":               {"test"},
		"property_url":                {"http://test.com"},
		"property_data":               {"test value"},
		"currency_identifier_desired": {"10"},
		"number_properties":           {"10"},
		"deadline":                    {"10"},
		"earlybird_bonus":             {"0.10"},
		"total_number":                {"10"},
	}

	return uri, result
}

func createCloseCrowdSaleTransaction() (string, url.Values) {
	uri := "/getunsigned/53"

	result := url.Values{
		"transaction_version": {"0"},
		"pubkey":              {""},
		"fee":                 {"0.01"},
		"transaction_from":    {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"currency_identifier": {"100"},
	}

	return uri, result
}

func createManagedIssuanceTransaction() (string, url.Values) {
	uri := "/getunsigned/54"

	result := url.Values{
		"transaction_version": {"0"},
		"pubkey":              {""},
		"fee":                 {"0.001"},
		"transaction_from":    {"bchtest:qrwejrccj3gfhtf025yt27fjhup09lxkfc0h0yl7ev"},
		"ecosystem":           {"1"},
		"precision":           {"1"},
		"previous_property_id": {"0"},
		"property_category":    {"shopping"},
		"property_subcategory": {"supermarket"},
		"property_name":        {"TSTM"},
		"property_url":         {"http://test.com"},
		"property_data":        {"test value"},
	}
	return uri, result
}

func createSendGrantTransaction() (string, url.Values) {
	uri := "/getunsigned/55"

	result := url.Values{
		"transaction_version": {"0"},
		"pubkey":              {""},
		"fee":                 {"0.001"},
		"transaction_from":    {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"currency_identifier": {"2"},
		"amount":              {"10.1"},
		"note":                {"test hello world"},
	}
	return uri, result
}

func createGetWhcTransaction() (string, url.Values) {
	uri := "/getunsigned/68"

	result := url.Values{
		"transaction_version": {"0"},
		"pubkey":              {""},
		"fee":                 {"0.01"},
		"transaction_from":    {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"amount_for_burn":     {"1"},
		"redeem_address":      {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
	}
	return uri, result
}

func createChangeIssuerTransaction() (string, url.Values) {
	uri := "/getunsigned/70"

	result := url.Values{
		"transaction_version": {"0"},
		"pubkey":              {""},
		"fee":                 {"0.001"},
		"transaction_from":    {"bchtest:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd"},
		"transaction_to":      {"bchtest:qrngs64rn42cana8nxqakyv0yw8ceg3sss0dfp4cp0"},
		"currency_identifier": {"1"},
	}
	return uri, result
}

func freezeTokenTransaction() (string, url.Values) {
	// rpc client instance
	uri := "/getunsigned/185"

	result := url.Values{
		"transaction_version": {"0"},
		"fee":                 {"0.0001"},
		"pubkey":              {""},
		"transaction_from":   {"bchtest:qrwejrccj3gfhtf025yt27fjhup09lxkfc0h0yl7ev"},
		"property_id":    {"459"},
		"amount":         {"1"},
		"frozen_address": {"bchtest:qqu9lh4jpc05p59pfhu9amyv9uvder8j3sa2up95vs"},
	}
	return uri, result
}

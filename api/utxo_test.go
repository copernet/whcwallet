package api

import (
	"strconv"
	"testing"

	"encoding/hex"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model/view"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

func TestGJson(t *testing.T) {
	content := []byte(`{"data":{"total_count":62,"page":1,"pagesize":50,
"list":[{"tx_hash":"17c40ce7f532445d2871e31e8276da961ab926d58ae4b255b1b69f4ed00bc42b","tx_output_n":7,"tx_output_n2":0,"value":2000696,"confirmations":267041},
{"tx_hash":"399eb8e47e558ebc7dde8004ab175edf06787fc5f5380aaef3d99cabc97694fc","tx_output_n":124,"tx_output_n2":0,"value":30911,"confirmations":267012},
{"tx_hash":"d5806457ffc74430839502bac6eb00b39871a41b39940ff4f582651f49bc77b4","tx_output_n":31,"tx_output_n2":0,"value":17135,"confirmations":267000},
{"tx_hash":"5adf096a85255d5e65e37417ac1ffab3752c2265584a9f3e4984cf624ba158e6","tx_output_n":1,"tx_output_n2":0,"value":3676290,"confirmations":266524},
{"tx_hash":"a329a603f299361f182d2e2e75b5e16e94840fe6f307dd04187954eda0fe2651","tx_output_n":1,"tx_output_n2":0,"value":545521,"confirmations":265525},
{"tx_hash":"ec21bfa9cc0ae1f0c60123e24f691e50fa9872e30126945b6d2d157552f6fc16","tx_output_n":0,"tx_output_n2":0,"value":2811341,"confirmations":261068},
{"tx_hash":"365e3141b7ae6c7612332b1ad0a657eccfdc4e029eb0167bda466e671a6874fe","tx_output_n":0,"tx_output_n2":0,"value":4210796,"confirmations":257109},
{"tx_hash":"0bb57f6e38012c86d4c5a28c904f2675082859147921a707d48961015a3e5057","tx_output_n":4,"tx_output_n2":0,"value":135906,"confirmations":253510},
{"tx_hash":"9facfdb6bb576e92ed273a983c2de06b7db5e683734d70e8e5dba5bb4ed0c0d0","tx_output_n":0,"tx_output_n2":0,"value":5403496,"confirmations":233516},
{"tx_hash":"6a650a26b2e6675a4b9afc0b24b2ef449f2e8334c794a1fe05d168f592990eee","tx_output_n":0,"tx_output_n2":0,"value":3179818,"confirmations":233363},
{"tx_hash":"7b8113b7e02ee7a4c17df0c4457275090e5fa208f506f2c912f42e2a5746b2ea","tx_output_n":0,"tx_output_n2":0,"value":3403554,"confirmations":233207},
{"tx_hash":"20eb5e9f70d53898b75a49ad955227fb8964a1f0a165140e373c71b80cc9efbf","tx_output_n":0,"tx_output_n2":0,"value":3260451,"confirmations":233057},
{"tx_hash":"7d910e16b5f71049ce47aa120c2e6d3346db4e37344559fa422d786a1fb628e9","tx_output_n":0,"tx_output_n2":0,"value":2854073,"confirmations":232605},
{"tx_hash":"9d59c3a4c8c3dd7f6d24118f07f8f169cdc8ae1c724c816bbf03570758c43331","tx_output_n":1,"tx_output_n2":0,"value":3133264,"confirmations":232455},
{"tx_hash":"0379571a44103b08915419c605ee30092ec057ec46ca12eb68efcb47c18ebad5","tx_output_n":0,"tx_output_n2":0,"value":2981560,"confirmations":232297},
{"tx_hash":"c3b7be7fb6cd9492ebbfb6bb47d8f4013cc57d531f403fc4ca8ef75b6e6b2d3c","tx_output_n":0,"tx_output_n2":0,"value":3069970,"confirmations":232130},
{"tx_hash":"9486126f0324e3818a4c328ac9771e57d281724f6a2c1b48e2e1bc20adf9eba8","tx_output_n":0,"tx_output_n2":0,"value":3299229,"confirmations":231972},
{"tx_hash":"f3c6cf12d387e2bac7e400c3ebfb5c92e733079f2028637cca70de7adc0e2099","tx_output_n":1,"tx_output_n2":0,"value":2400116,"confirmations":231437},
{"tx_hash":"4ef49146f251c727ac9e024ce2c8a41afad4b6a709d644ae378514cbce91e0e8","tx_output_n":0,"tx_output_n2":0,"value":2618114,"confirmations":231324},
{"tx_hash":"a657c3a56d975a9ced1f2b26dfea638ae3f64545dc8e741cbf094c688aaee136","tx_output_n":1,"tx_output_n2":0,"value":3034977,"confirmations":231175},
{"tx_hash":"85bd4c45e162ad3d047a46631669f92d6798437000ea58477aeaa99bf3754be2","tx_output_n":0,"tx_output_n2":0,"value":3108415,"confirmations":231041},
{"tx_hash":"8861a0a9b14be5d895b5774f33dd295e8cba0cff2aad2da16ffbea6e3d82697a","tx_output_n":0,"tx_output_n2":0,"value":1584577,"confirmations":230874}]},
"err_no":0,"err_msg":null}`)

	contentStr := string(content)
	if gjson.Get(contentStr, "err_no").String() != "0" {
		t.Error("parse err_no error")
	}
	if gjson.Get(contentStr, "data.list.#").Int() != 22 {
		t.Error("get the number of elements error")
	}

	getList := gjson.Get(contentStr, "data.list").Array()
	if len(getList) != 22 {
		t.Error("get array length error")
	}

	getList[5].ForEach(func(key, value gjson.Result) bool {
		if key.String() == "tx_hash" {
			if value.String() != "ec21bfa9cc0ae1f0c60123e24f691e50fa9872e30126945b6d2d157552f6fc16" {
				t.Error("itera array error in foreach")
			}
		}
		return true
	})
}

func TestGetUtxo(t *testing.T) {
	address := "bchtest:qpz8py2yqyp7x2aaqfvrwdk4jf2c23ypzczr6weclj"
	addr, err := cashutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	if err != nil {
		t.Error("decode bitcoin address error")
	}

	ret, err := getUtxo(addr, decimal.NewFromFloat(0.0000001), 1, nil, 0)
	if err != nil {
		t.Error(err.Error())
	}

	// the result depends the specified address
	if len(ret) < 1 {
		t.Error("empty result")
	}
	t.Log("get unspent list length:" + strconv.Itoa(len(ret)))
}

func TestGetUtxoElectrumx(t *testing.T) {
	required := decimal.NewFromFloat(0.0002)
	addr, err := cashutil.DecodeAddress("bchtest:qrpcfpwt8ptjpj4j5r2rn69s4x0psfwmxqddjqtlq7", &chaincfg.TestNet3Params)
	if err != nil {
		t.Error(err.Error())
	}

	ret, err := GetUtxoElectrumx(addr, required, false)
	if err != nil {
		t.Error(err.Error())
	}
	if len(ret) < 1 {
		t.Error("empty result")
	}

	t.Log("get unspent list length:" + strconv.Itoa(len(ret)))

}

func TestAssembleBCHTx(t *testing.T) {
	addr1 := "moegtR1zdP4ZFFQ9UDuepZtfxL2tcBvbXt"
	addr2 := "2NBD5C5VWaDzdt6ULxrBk3vpETHCbdQSASz"
	b1, _ := getScriptPubByAddr(addr1)
	t.Log(hex.EncodeToString(b1))
	b2, _ := getScriptPubByAddr(addr2)
	t.Log(hex.EncodeToString(b2))
}

func TestDoAssemble(t *testing.T) {
	tx := view.UserTx{
		FromAddress: "moegtR1zdP4ZFFQ9UDuepZtfxL2tcBvbXt",
		ToAddress: []string{
			"2NBD5C5VWaDzdt6ULxrBk3vpETHCbdQSASz",
			"mrefpMHh5x7KyFvBSFnt37Bq38pnMV84N6",
		},
		ToAmount: []float64{
			0.01,
			0.02,
		},
		Fee: 0.00003,
	}

	tx2 := view.UserTx{
		FromAddress: "bchtest:qpvnw78dtuglcl4l6sv9rgysz39eywd9vgdfzggac0",
		ToAddress: []string{
			"bchtest:przs59k0lqekuf72u5cp8j69nwzhsmr5lgysnq8m0j",
			"bchtest:qpapaaxm8u3prc3k9xqmx0r4q4v230ewd58yzc0chp",
		},
		ToAmount: []float64{
			0.01,
			0.02,
		},
		Fee: 0.00003,
	}
	hex, utxoes, err := doTheBCHAssemble(&tx, &gin.Context{})
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(*hex)
		t.Log(utxoes)
	}

	hex2, utxo2es, err := doTheBCHAssemble(&tx2, &gin.Context{})
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(*hex2)
		t.Log(utxo2es)
	}

}

func TestAddrConv(t *testing.T) {
	var addrs = []string{"moegtR1zdP4ZFFQ9UDuepZtfxL2tcBvbXt", "2NBD5C5VWaDzdt6ULxrBk3vpETHCbdQSASz", "mrefpMHh5x7KyFvBSFnt37Bq38pnMV84N6"}
	for _, value := range addrs {
		addr, _ := cashutil.DecodeAddress(value, config.GetChainParam())
		t.Log(value + ":" + addr.EncodeAddress(true))
	}
}

func TestIsAvailableCoin(t *testing.T) {
	tests := []struct {
		hash  string
		index int
		ok    bool
	}{
		{
			hash:  "8326bf7f97f86b012e64d26abb9019dd63ab3ae7d023d5a6bb1a3a157ed29ed2",
			index: 1,
			ok:    false,
		},
		{
			hash:  "5e03b49aad01b172195afacae92ce38daa930e27a9ef054012ad518bc53f5845",
			index: 2,
			ok:    true,
		},
	}

	for _, test := range tests {
		ok, _, err := isAvailableCoin(test.hash, test.index)
		if test.ok && err != nil {
			t.Errorf("bad result, the expected result is %v, but got %s. the transaciont hash: %s,",
				test.ok, err.Error(), test.hash)
		}

		if ok != test.ok {
			t.Errorf("bad check, the expecte result is %v, but got %v. the transaction hash: %s",
				test.ok, ok, test.hash)
		}
	}
}

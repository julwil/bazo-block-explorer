var app = new Vue({
    el: '#createTxModal',
    methods: {
        saveTransaction: function () {
            const headers = {
                'Content-Type': 'text/plain'
            };
            axios.post(`${app.baseUrl}/createFundsTx`, {
                from: this.funds.from,
                to: this.funds.to,
                amount: +this.funds.amount,
                txcount: +this.funds.txCount,
                chparams: this.chParams,
                data: this.data,
            }, {headers}).then((response) => {
                $
                const responseBody = response.data || {};
                const data = responseBody.content[0];
                const txHash = data.detail;

                const curve = new elliptic.ec('p256');
                const privateKey = curve.keyFromPrivate(this.privateKey);
                const signature = privateKey.sign(txHash);

                axios.post(`${app.baseUrl}/signFundsTx`, {
                    hash: txHash,
                    signature: signature.r.toJSON() + signature.s.toJSON()
                }, {headers}).then((response) => {
                    let elementId = '';
                    response.status >= 300 ? elementId = 'alert-fail' : 'alert-success';
                    $(`#${elementId}`).toggle();
                })
            });
        }
    },
    data: {
        baseUrl: 'http://localhost:8010',
        txType: '',
        publicKey: '',
        privateKey: '',
        chParams: '',
        funds: {
            from: '',
            to: '',
            txCount: '',
            amount: '',
        },
        data: ''
    },
    delimiters: ["<%", "%>"]
})

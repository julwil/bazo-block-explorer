var app = new Vue({
    el: '#transactions-vue-container',
    methods: {
        createAccountTx: function () {
            const headers = {
                'Content-Type': 'text/plain'
            };
            axios.post(`${this.baseUrl}/tx/acc`, {
                root_wallet: this.account.rootWallet,
                wallet: this.account.wallet,
                ch_params: this.chParams,
                fee: parseInt(this.fee),
                data: this.data,
            }, {headers}).then((response) => this.createTxResponseHandler(response));
        },

        createFundsTx: function () {
            const headers = {
                'Content-Type': 'text/plain'
            };
            axios.post(`${this.baseUrl}/tx/funds`, {
                from: this.funds.from,
                to: this.funds.to,
                amount: parseInt(this.funds.amount),
                tx_count: parseInt(this.funds.txCount),
                fee: parseInt(this.fee),
                ch_params: this.chParams,
                data: this.data,
            }, {headers}).then((response) => this.createTxResponseHandler(response));
        },

        createUpdateTx: function () {
            const headers = {
                'Content-Type': 'text/plain'
            };
            axios.post(`${this.baseUrl}/tx/update`, {
                tx_to_update: this.update.txToUpdate,
                tx_issuer: this.update.txToUpdateIssuer,
                update_data: this.update.updateData,
                ch_params: this.chParams,
                fee: parseInt(this.fee),
                data: this.data,
            }, {headers}).then((response) => this.createTxResponseHandler(response));
        },

        createTxHandler: function () {
            switch (this.txType) {
                case 'Account Tx':
                    return this.createAccountTx();
                case 'Funds Tx':
                    return this.createFundsTx();
                case 'Update Tx':
                    return this.createUpdateTx();
                default:
                    return;
            }
        },

        signTxResponseHandler: function (response) {
            const elementId = response.status >= 300 ? 'alert-fail' : 'alert-success';
            $("#" + elementId).toggle();
        },

        createTxResponseHandler: function (response) {
            const responseBody = response.data || {};
            const data = responseBody.content[0];
            this.txHash = data.detail;

            $("#createTxModal").modal('hide');
            $("#signTxModal").modal('show');
        },

        signTxHash: function (txHash, privateKeyString) {
            const curve = new elliptic.ec('p256');
            const privateKey = curve.keyFromPrivate(privateKeyString);
            const signature = privateKey.sign(txHash);

            return signature.r.toJSON() + signature.s.toJSON();
        },

        signTxHandler: function() {
            if (this.privateKey.length === 0) {
                $("#alert-failed").show();
                return;
            }

            const headers = {
                'Content-Type': 'text/plain'
            };

            const signature = this.signTxHash(this.txHash, this.privateKey);
            axios.post(`${this.baseUrl}/tx/signature`, {
                hash: this.txHash,
                signature: signature
            }, {headers})
                .then((response) => this.signTxResponseHandler(response))

            setTimeout(() => location.reload(), 3000);
        }
    },
    data: {
        baseUrl: 'http://api.bazo.local',
        txType: '',
        txHash: '',
        privateKey: '',
        chParams: '',
        fee: 1,
        data: '',
        funds: {
            from: '',
            to: '',
            txCount: '',
            amount: '',
        },
        account: {
            rootWallet: '',
            wallet: '',
        },
        update: {
            txToUpdate: '',
            txToUpdateIssuer: '',
            updateData: '',
        }
    },
    delimiters: ["<%", "%>"]
})

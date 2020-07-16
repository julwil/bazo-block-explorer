var app = new Vue({
  el: '#parameter-list',
  methods:{
    isPositiveInt: function(value) {
      return value >>> 0 === parseFloat(value);
    },
    closeModal: function () {
      $('#myModal').modal('hide');
    },
    setPublicKeyFromCookie: function () {
      app.accountinfo.rootpublickey = this.$cookies.get("publicKey")
      console.log(app.accountinfo.rootpublickey)
      console.log(this.$cookies.get("publicKey"))
    },
    signTransaction: function (txhash, privatekey) {
      console.log(txhash, privatekey);
      var curve = new elliptic.ec('p256')
      var privatekey = curve.keyFromPrivate(app.accountinfo.privatekey)
      var signature = privatekey.sign(app.accountinfo.txhash)
      var signatureHexString = signature.r.toJSON() + signature.s.toJSON()

      axios.post(`${app.baseUrl}/sendConfigTx/${app.accountinfo.txhash}/${signatureHexString}`).then(
        function() {
          app.closeModal()
        }
      )
    },
    changeBlockSize: function(blocksize, fee) {
      console.log("Changing blocksize: ", blocksize, fee)
      app.setPublicKeyFromCookie()
      if (app.isPositiveInt(blocksize) && app.isPositiveInt(fee)) {
        axios.get(`${app.baseUrl}/account/'${app.accountinfo.rootpublickey}`).then(
          function(response) {
            console.log(response.data)
            app.accountinfo.txcount = response.data.txCnt
            if (response.data.content[0].detail.isRoot) {
              axios.post(`${app.baseUrl}/createConfigTx/${0}/${1}/${app.blocksize.blocksize}/${app.blocksize.fee}/${app.accountinfo.txcount}`).then(
                function(response) {
                  $("#transactionModal").modal()
                  console.log(response.data);
                  app.accountinfo.txhash = response.data.content[0].detail
                }
              )
            }
          }
        )
      } else {
        alert("Values are not positive integers!");
      }
    },
    changeDiffInterval: function(diffinterval, fee) {
      console.log("Changing difficulty interval: ", diffinterval, fee)
      app.setPublicKeyFromCookie()
      if (app.isPositiveInt(diffinterval) && app.isPositiveInt(fee)) {
        axios.get(`${app.baseUrl}/account/${app.accountinfo.rootpublickey}`).then(
          function(response) {
            console.log(response.data)
            app.accountinfo.txcount = response.data.txCnt
            if (response.data.isRoot) {
              axios.post(`${app.baseUrl}/createConfigTx/${0}/${2}/${app.diffinterval.diffinterval}/${app.diffinterval.fee}/${app.accountinfo.txcount}`).then(
                function(response) {
                  $("#myModal").modal()
                  console.log(response.data);
                  app.accountinfo.txhash = response.data.content[0].detail
                }
              )
            }
          }
        )
      } else {
        alert("Values are not positive integers!");
      }
    },
    changeMinFee: function(minfee, fee) {
      console.log("Changing minimum fee: ", minfee, fee)
      app.setPublicKeyFromCookie()
      if (app.isPositiveInt(minfee) && app.isPositiveInt(fee)) {
        axios.get(`${app.baseUrl}/account/${app.accountinfo.rootpublickey}`).then(
          function(response) {
            console.log(response.data)
            app.accountinfo.txcount = response.data.txCnt
            if (response.data.isRoot) {
              axios.post(`${app.baseUrl}/createConfigTx/${0}/${3}/${app.minfee.minfee}/${app.minfee.fee}/${app.accountinfo.txcount}`).then(
                function(response) {
                  $("#myModal").modal()
                  console.log(response.data);
                  app.accountinfo.txhash = response.data.content[0].detail
                }
              )
            }
          }
        )
      } else {
        alert("Values are not positive integers!");
      }
    },
    changeBlockInterval: function(blockinterval, fee) {
      console.log("Changing block interval: ", blockinterval, fee)
      app.setPublicKeyFromCookie()
      if (app.isPositiveInt(blockinterval) && app.isPositiveInt(fee)) {
        axios.get(`${app.baseUrl}/account/${app.accountinfo.rootpublickey}`).then(
          function(response) {
            console.log(response.data)
            app.accountinfo.txcount = response.data.txCnt
            if (response.data.isRoot) {
              axios.post(`${app.baseUrl}/createConfigTx/${0}/${4}/${app.blockinterval.blockinterval}/${app.blockinterval.fee}/${app.accountinfo.txcount}`).then(
                function(response) {
                  $("#myModal").modal()
                  console.log(response.data);
                  app.accountinfo.txhash = response.data.content[0].detail
                }
              )
            }
          }
        )
      } else {
        alert("Values are not positive integers!");
      }
    },
    changeBlockReward: function(blockreward, fee) {
      console.log("Changing block reward: ", blockreward, fee)
      app.setPublicKeyFromCookie()
      if (app.isPositiveInt(blockreward) && app.isPositiveInt(fee)) {
        axios.get(`${app.baseUrl}/account/${app.accountinfo.rootpublickey}`).then(
          function(response) {
            console.log(response.data)
            app.accountinfo.txcount = response.data.txCnt
            if (response.data.isRoot) {
              axios.post(`${app.baseUrl}/createConfigTx/${0}/${5}/${app.blockreward.blockreward}/${app.blockreward.fee}/${app.accountinfo.txcount}`).then(
                function(response) {
                  $("#myModal").modal()
                  console.log(response.data);
                  app.accountinfo.txhash = response.data.content[0].detail
                }
              )
            }
          }
        )
      } else {
        alert("Not all values are positive integers!");
      }
    },
  },
  data: {
    baseUrl: 'http://localhost:8010',
    blocksize: {
      blocksize: '',
      fee: ''
    },
    diffinterval: {
      diffinterval: '',
      fee: ''
    },
    minfee: {
      minfee: '',
      fee: ''
    },
    blockinterval: {
      blockinterval: '',
      fee: ''
    },
    blockreward: {
      blockreward: '',
      fee: ''
    },
    accountinfo: {
      privatekey: '',
      rootpublickey: '',
      txhash: '',
      txcount: ''
    },

    tx: {
        type: '',
    }
  },
  delimiters: ["<%","%>"]
})


var app = new Vue({
  el: '#parameter-list',
  methods:{
    closeModal: function () {
      $('#myModal').modal('hide');
    },
    signTransaction: function (txhash, privatekey) {
      console.log(txhash, privatekey);
      var curve = new elliptic.ec('p256')
      var privatekey = curve.keyFromPrivate(app.accountinfo.privatekey)
      var signature = privatekey.sign(app.accountinfo.txhash)
      var signatureHexString = signature.r.toJSON() + signature.s.toJSON()

      axios.post(`http://192.41.136.199:8001/sendConfigTx/${app.accountinfo.txhash}/${signatureHexString}`).then(
        function() {
          app.closeModal()
        }
      )
    },
    changeBlockSize: function(blocksize, fee) {
      console.log("Changing blocksize: ", blocksize, fee)

      axios.get('http://192.41.136.199:8001/account/' + app.accountinfo.rootpublickey).then(
        function(response) {
          console.log(response.data)
          app.accountinfo.txcount = response.data.txCnt
          if (response.data.isRoot) {
            axios.post(`http://192.41.136.199:8001/createConfigTx/${0}/${1}/${blocksize.blocksize}/${app.blocksize.fee}/${app.accountinfo.txcount}`).then(
              function(response) {
                $("#myModal").modal()
                console.log(response.data);
                app.accountinfo.txhash = response.data.content[0].detail
              }
            )
          }
        }
      )

    }
  },
  data: {
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
      rootpublickey: 'f894ba7a24c1c324bc4b0a833d4b076a0e0f675a380fb7e782672c6568aaab0669ddbc62f79cb521411840d83ff0abf941a8e717d81af3dfc2973f1bac30308a',
      txhash: '',
      txcount: ''
    },

    title: 'Change System Parameters'
  },
  delimiters: ["<%","%>"]
})

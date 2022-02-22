export default {
  tx: {
    txid: '5E4DE0462EA7F140C122F12B92BEF09C55F430079FABBF426710847006EF1935',
    blockHeight: '9473608',
    fee: {
      amount: '2500',
      denom: 'uatom',
    },
    gasUsed: '62326',
    gasWanted: '85000',
    index: 3,
    value: '',
    messages: [
      {
        from: 'cosmos1t5u0jfg3ljsjrh2m9e47d4ny2hea7eehxrzdgd',
        to: 'cosmos14e25lpsedq863vgweqg4m9n0z28c203kfdlzmz',
        type: 'send',
        value: {
          amount: '2002965',
          denom: 'uatom',
        },
      },
    ],
    events: [
      {
        type: 'coin_received',
        attributes: [
          {
            key: 'receiver',
            value: 'cosmos14e25lpsedq863vgweqg4m9n0z28c203kfdlzmz',
          },
          {
            key: 'amount',
            value: '2002965uatom',
          },
        ],
      },
      {
        type: 'coin_spent',
        attributes: [
          {
            key: 'spender',
            value: 'cosmos1t5u0jfg3ljsjrh2m9e47d4ny2hea7eehxrzdgd',
          },
          {
            key: 'amount',
            value: '2002965uatom',
          },
        ],
      },
      {
        type: 'message',
        attributes: [
          {
            key: 'action',
            value: '/cosmos.bank.v1beta1.MsgSend',
          },
          {
            key: 'sender',
            value: 'cosmos1t5u0jfg3ljsjrh2m9e47d4ny2hea7eehxrzdgd',
          },
          {
            key: 'module',
            value: 'bank',
          },
        ],
      },
      {
        type: 'transfer',
        attributes: [
          {
            key: 'recipient',
            value: 'cosmos14e25lpsedq863vgweqg4m9n0z28c203kfdlzmz',
          },
          {
            key: 'sender',
            value: 'cosmos1t5u0jfg3ljsjrh2m9e47d4ny2hea7eehxrzdgd',
          },
          {
            key: 'amount',
            value: '2002965uatom',
          },
        ],
      },
    ],
  },
}

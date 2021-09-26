import * as yup from 'yup';

const addRendezvousSchema = yup.object({
  rendezvousAddress: yup
    .string()
    .defined('Address is required')
    .test('Address test', 'Invalid address', function (value) {
      let regex = /\/ip4\/\d+\.\d+\.\d+\.\d+\/tcp\/\d+\/p2p\/\w+/gm;
      if (value == null || value == '' || value == undefined) {
        return false;
      }

      let matches = value.match(regex);

      return matches && matches.length > 0;
    })
});

export default addRendezvousSchema;

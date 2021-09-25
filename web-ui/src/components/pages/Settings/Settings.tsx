import { Box, TextField } from '@material-ui/core';
import AddRoundedIcon from '@material-ui/icons/AddRounded';
import { useFormik } from 'formik';
import { FC, useState } from 'react';
import RendezvousService from '../../../services/rendezvous/rendezvousService';
import addRendezvousSchema from '../../../shared/schemas/rendezvousSchemas';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import FormTitle from '../../atoms/FormTitle/FormTitle';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import RendezvousList from '../../molecules/RendezvousList/RendezvousList';
import useSnackbar from '../../molecules/Snackbar/useSnackbar.hook';
import { ISettingsProps } from './settings.types';

const Settings: FC<ISettingsProps> = () => {
  const { openSnackbar } = useSnackbar();

  const [trigger, setTrigger] = useState<boolean>(false);

  const formik = useFormik({
    initialValues: {
      rendezvousAddress: ''
    },
    enableReinitialize: true,
    validationSchema: addRendezvousSchema,
    onSubmit: (values, { resetForm }) => {
      const handleCreate = async () => {
        return await RendezvousService.addNewRendezvousNode({
          address: values.rendezvousAddress
        });
      };

      handleCreate()
        .then((response) => {
          setTrigger(!trigger);
          openSnackbar('Rendezvous successfully added', 'success');
        })
        .catch((err) => {
          openSnackbar('Unable to add rendezvous', 'error');
        });

      resetForm();
    }
  });

  return (
    <Box
      display={'flex'}
      flexDirection={'column'}
      width={'80%'}
      height={'100%'}
    >
      <Box display={'flex'} width={'100%'} mb={4}>
        <PageTitle title={'Settings'} />
      </Box>
      <Box display={'flex'} width={'100%'} flexDirection={'column'}>
        <form onSubmit={formik.handleSubmit} autoComplete={'off'}>
          <Box mb={2}>
            <FormTitle title={'Add a rendezvous node'} />
          </Box>
          <Box display={'flex'} alignItems={'center'}>
            <Box height={'80px'}>
              <TextField
                id={'rendezvousAddress'}
                label={'Rendezvous address'}
                variant={'outlined'}
                value={formik.values.rendezvousAddress}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                placeholder={'/ip4/x.x.x.x/tcp/port/p2p/Qm...'}
                error={
                  formik.touched.rendezvousAddress &&
                  Boolean(formik.errors.rendezvousAddress)
                }
                helperText={
                  formik.touched.rendezvousAddress &&
                  formik.errors.rendezvousAddress
                }
              />
            </Box>
            <Box alignItems={'center'} mb={'24px'} ml={2}>
              <ActionButton
                text={'Add rendezvous'}
                startIcon={<AddRoundedIcon />}
              />
            </Box>
          </Box>
        </form>

        <Box mt={4} width={'100%'}>
          <RendezvousList trigger={trigger} setTrigger={setTrigger} />
        </Box>
      </Box>
    </Box>
  );
};

export default Settings;

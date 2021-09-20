import {
  Backdrop,
  Box,
  Fade,
  IconButton,
  Modal,
  TextField,
  Typography
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import CloseRoundedIcon from '@material-ui/icons/CloseRounded';
import { useFormik } from 'formik';
import { FC } from 'react';
import { joinWorkspacePasswordSchema } from '../../../shared/schemas/workspaceSchemas';
import ModalButton from '../../atoms/ModalButton/ModalButton';
import { EModalButtonType } from '../../atoms/ModalButton/modalButton.types';
import { IJoinWorkspacePasswordModalProps } from './joinWorkspacePasswordModal.types';

const JoinWorkspacePasswordModal: FC<IJoinWorkspacePasswordModalProps> = (
  props
) => {
  const { open, handleConfirm } = props;
  const classes = useStyles();

  const formik = useFormik({
    initialValues: {
      password: ''
    },
    enableReinitialize: true,
    validationSchema: joinWorkspacePasswordSchema,
    onSubmit: (values, { resetForm }) => {
      handleConfirm(values.password, true);

      resetForm();
    }
  });

  return (
    <Modal
      className={classes.modal}
      open={open}
      onClose={() => handleConfirm('', false)}
      closeAfterTransition
      BackdropComponent={Backdrop}
      BackdropProps={{
        timeout: 500
      }}
    >
      <Fade in={open}>
        <div className={classes.modalWrapper}>
          <form autoComplete={'off'} onSubmit={formik.handleSubmit}>
            <Box
              display={'flex'}
              alignItems={'center'}
              justifyContent={'space-between'}
            >
              <Typography className={classes.modalTitle}>
                Enter workspace password
              </Typography>
              <IconButton
                classes={{
                  root: 'iconButtonRoot'
                }}
                onClick={() => handleConfirm('', false)}
              >
                <CloseRoundedIcon
                  style={{
                    width: '20px',
                    height: 'auto'
                  }}
                />
              </IconButton>
            </Box>
            <Box
              display={'flex'}
              flexDirection={'column'}
              mt={5}
              width={'100%'}
            >
              <Typography className={classes.modalTextMain}>
                Workspace password:
              </Typography>
              <Box width={'100%'} mt={2} minHeight={'80px'}>
                <TextField
                  id={'password'}
                  type={'password'}
                  label={'Password'}
                  variant={'outlined'}
                  value={formik.values.password}
                  onChange={formik.handleChange}
                  onBlur={formik.handleBlur}
                  error={
                    formik.touched.password && Boolean(formik.errors.password)
                  }
                  helperText={formik.touched.password && formik.errors.password}
                />
              </Box>
            </Box>
            <Box
              display={'flex'}
              alignItems={'center'}
              justifyContent={'center'}
              width={'100%'}
              mt={5}
            >
              <ModalButton
                handleConfirm={() => {
                  formik.submitForm();
                }}
                type={EModalButtonType.FILLED}
                text={'Confirm'}
              />
              <ModalButton
                handleConfirm={() => handleConfirm('', false)}
                type={EModalButtonType.OUTLINED}
                text={'Cancel'}
                margins={'0 0 0 16px'}
              />
            </Box>
          </form>
        </div>
      </Fade>
    </Modal>
  );
};

const useStyles = makeStyles((theme) => {
  return {
    modal: {
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center'
    },
    modalTitle: {
      fontWeight: 600,
      fontSize: '1.2rem',
      textAlign: 'left'
    },
    modalWrapper: {
      backgroundColor: theme.palette.background.paper,
      boxShadow: theme.palette.boxShadows.darker,
      padding: '20px 30px',
      // TODO add responsive features
      width: '400px',
      height: 'auto',
      border: 'none',
      outline: 'none',
      borderRadius: '15px'
    },
    modalTextMain: {},
    modalSubtext: {
      fontWeight: 500
    }
  };
});

export default JoinWorkspacePasswordModal;

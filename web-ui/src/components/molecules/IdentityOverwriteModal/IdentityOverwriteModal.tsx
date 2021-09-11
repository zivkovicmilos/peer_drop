import {
  Backdrop,
  Box,
  Fade,
  IconButton,
  Modal,
  Typography
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import CloseRoundedIcon from '@material-ui/icons/CloseRounded';
import { FC } from 'react';
import ModalButton from '../../atoms/ModalButton/ModalButton';
import { EModalButtonType } from '../../atoms/ModalButton/modalButton.types';
import { IIdentityOverwriteModalProps } from './identityOverwriteModal.types';

const IdentityOverwriteModal: FC<IIdentityOverwriteModalProps> = (props) => {
  const { open, publicKeyID, handleConfirm } = props;
  const classes = useStyles();

  const renderKeyInfo = () => {
    return (
      <Typography
        className={classes.modalSubtext}
      >{`Public key ID · ${publicKeyID}`}</Typography>
    );
  };

  return (
    <Modal
      className={classes.modal}
      open={open}
      onClose={() => handleConfirm(false)}
      closeAfterTransition
      BackdropComponent={Backdrop}
      BackdropProps={{
        timeout: 500
      }}
    >
      <Fade in={open}>
        <div className={classes.modalWrapper}>
          <Box
            display={'flex'}
            alignItems={'center'}
            justifyContent={'space-between'}
          >
            <Typography className={classes.modalTitle}>
              Are you sure?
            </Typography>
            <IconButton onClick={() => handleConfirm(false)}>
              <CloseRoundedIcon
                style={{
                  width: '20px',
                  height: 'auto'
                }}
              />
            </IconButton>
          </Box>
          <Box display={'flex'} flexDirection={'column'} mt={5}>
            <Typography className={classes.modalTextMain}>
              You are about to overwrite:
            </Typography>
            {renderKeyInfo()}
          </Box>
          <Box
            display={'flex'}
            alignItems={'center'}
            justifyContent={'center'}
            width={'100%'}
            mt={5}
          >
            <ModalButton
              handleConfirm={() => handleConfirm(true)}
              type={EModalButtonType.FILLED}
              text={'Confirm'}
            />
            <ModalButton
              handleConfirm={() => handleConfirm(false)}
              type={EModalButtonType.OUTLINED}
              text={'Cancel'}
              margins={'0 0 0 16px'}
            />
          </Box>
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
      fontSize: '1.5rem'
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

export default IdentityOverwriteModal;

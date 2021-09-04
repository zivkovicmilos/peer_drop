import { Backdrop, Fade, Modal } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { FC } from 'react';
import { IContactsConfirmModalProps } from './contactsConfirmModal.types';

const ContactsConfirmModal: FC<IContactsConfirmModalProps> = (props) => {
  const { open, contactInfo, handleConfirm } = props;
  const classes = useStyles();

  console.log(contactInfo);

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
        <div className={classes.paper}>
          <h2 id="transition-modal-title">Transition modal</h2>
          <p id="transition-modal-description">
            react-transition-group animates me.
          </p>
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
    paper: {
      backgroundColor: theme.palette.background.paper,
      border: '2px solid #000',
      boxShadow: theme.shadows[5],
      padding: theme.spacing(2, 4, 3)
    }
  };
});

export default ContactsConfirmModal;

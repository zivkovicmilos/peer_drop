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
import {
  ENewWorkspaceType,
  ENWAccessControl
} from '../../../context/newWorkspaceContext.types';
import { IWorkspaceInfoResponse } from '../../../services/workspaces/workspacesService.types';
import { ReactComponent as JoinImage } from '../../../shared/assets/icons/undraw_join_of2w.svg';
import { IJoinWorkspaceModalProps } from './joinWorkspaceModal.types';

const JoinWorkspaceModal: FC<IJoinWorkspaceModalProps> = (props) => {
  const { open, workspaceInfo, handleConfirm } = props;
  const classes = useStyles();

  const renderAdditionalSettings = () => {
    let additionalSettings = [];

    // if (workspaceInfo.permissions.additionalOwners.active) {
    //   additionalSettings.push(
    //     <Typography key={'additional-owners'} className={classes.reviewItem}>
    //       {`Additional workspace owners - ${workspaceInfo.permissions.additionalOwners.contactIDs.length}`}
    //     </Typography>
    //   );
    // }
    //
    // if (workspaceInfo.permissions.autocloseWorkspace.active) {
    //   additionalSettings.push(
    //     <Typography key={'autoclose'} className={classes.reviewItem}>
    //       {`Auto-close workspace - ${workspaceInfo.permissions.autocloseWorkspace.date}`}
    //     </Typography>
    //   );
    // }
    //
    // if (workspaceInfo.permissions.enforcePeerLimit.active) {
    //   additionalSettings.push(
    //     <Typography key={'peer-limit'} className={classes.reviewItem}>
    //       {`Enforce peer limit - ${workspaceInfo.permissions.enforcePeerLimit.limit}`}
    //     </Typography>
    //   );
    // }

    if (additionalSettings.length < 1) {
      additionalSettings.push(
        <Typography key={'no-settings'} className={classes.reviewItem}>
          No additional settings
        </Typography>
      );
    }

    return additionalSettings;
  };

  const convertWorkspaceType = (type: string): ENewWorkspaceType => {
    switch (type) {
      case 'send-only':
        return ENewWorkspaceType.SEND_ONLY;
      case 'receive-only':
        return ENewWorkspaceType.RECEIVE_ONLY;
      default:
        return ENewWorkspaceType.SEND_RECEIVE;
    }
  };

  const convertAccessControlType = (type: string): ENWAccessControl => {
    if (type == 'password') {
      return ENWAccessControl.PASSWORD;
    } else {
      return ENWAccessControl.SPECIFIC_CONTACTS;
    }
  };

  const convertAccessControl = (
    workspaceInfo: IWorkspaceInfoResponse
  ): { contacts: string[] } | { password: string } => {
    if (workspaceInfo.passwordHash) {
      return {
        password: workspaceInfo.passwordHash
      };
    } else {
      return {
        contacts: workspaceInfo.contactsWrapper.contactPublicKeys
      };
    }
  };

  const renderBaseSettings = () => {
    let baseSettings = [];

    switch (workspaceInfo.workspaceType) {
      case ENewWorkspaceType.SEND_ONLY: {
        baseSettings.push(
          <Typography key={'send-only'} className={classes.reviewItem}>
            {ENewWorkspaceType.SEND_ONLY}
          </Typography>
        );
        break;
      }
      case ENewWorkspaceType.RECEIVE_ONLY: {
        baseSettings.push(
          <Typography key={'receive-only'} className={classes.reviewItem}>
            {ENewWorkspaceType.RECEIVE_ONLY}
          </Typography>
        );
        break;
      }
      case ENewWorkspaceType.SEND_RECEIVE: {
        baseSettings.push(
          <Typography key={'send-receive'} className={classes.reviewItem}>
            {ENewWorkspaceType.SEND_RECEIVE}
          </Typography>
        );
        break;
      }
    }

    if (
      convertAccessControlType(workspaceInfo.securityType) ==
      ENWAccessControl.PASSWORD
    ) {
      baseSettings.push(
        <Typography key={'access-password'} className={classes.reviewItem}>
          {ENWAccessControl.PASSWORD}
        </Typography>
      );
    } else {
      baseSettings.push(
        <Typography key={'access-contacts'} className={classes.reviewItem}>
          {ENWAccessControl.SPECIFIC_CONTACTS}
        </Typography>
      );
    }

    return baseSettings;
  };

  if (workspaceInfo) {
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
                Joined workspace!
              </Typography>
              <IconButton
                classes={{
                  root: 'iconButtonRoot'
                }}
                onClick={() => handleConfirm(false)}
              >
                <CloseRoundedIcon
                  style={{
                    width: '20px',
                    height: 'auto'
                  }}
                />
              </IconButton>
            </Box>
            <Box display={'flex'} mt={2} width={'100%'}>
              <Box display={'flex'} flexDirection={'column'}>
                <Typography>You've joined the workspace:</Typography>
                <Box my={1}>
                  <Typography className={classes.workspaceName}>
                    {workspaceInfo.name}
                  </Typography>
                </Box>
                {renderBaseSettings()}
                {renderAdditionalSettings()}
              </Box>
              <Box ml={'auto'}>
                <JoinImage
                  style={{
                    width: '250px',
                    height: 'auto'
                  }}
                />
              </Box>
            </Box>
          </div>
        </Fade>
      </Modal>
    );
  } else {
    return null;
  }
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
      width: '600px',
      height: 'auto',
      border: 'none',
      outline: 'none',
      borderRadius: '15px'
    },
    workspaceName: {
      fontWeight: 600,
      fontSize: '1.2rem'
    },
    modalTextMain: {},
    modalSubtext: {
      fontWeight: 500
    },
    reviewItem: {
      fontWeight: 600,
      color: '#9C9C9C',
      fontSize: theme.typography.pxToRem(16)
    }
  };
});

export default JoinWorkspaceModal;

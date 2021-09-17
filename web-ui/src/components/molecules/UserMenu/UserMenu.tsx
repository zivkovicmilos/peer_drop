import {
  Avatar,
  Box,
  Button,
  ClickAwayListener,
  Grow,
  MenuItem,
  MenuList,
  Paper,
  Popper,
  Typography
} from '@material-ui/core';
import { makeStyles, Theme } from '@material-ui/core/styles';
import KeyboardArrowDownRoundedIcon from '@material-ui/icons/KeyboardArrowDownRounded';
import KeyboardArrowUpRoundedIcon from '@material-ui/icons/KeyboardArrowUpRounded';
import PowerSettingsNewRoundedIcon from '@material-ui/icons/PowerSettingsNewRounded';
import SwapHorizRoundedIcon from '@material-ui/icons/SwapHorizRounded';
import React, { FC, Fragment, useContext, useRef, useState } from 'react';
import SessionContext from '../../../context/SessionContext';
import CommonService from '../../../services/common/commonService';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import SwapIdentities from '../SwapIdentities/SwapIdentities';
import { IUserMenuProps } from './userMenu.types';

const UserMenu: FC<IUserMenuProps> = () => {
  const { userIdentity } = useContext(SessionContext);
  const [open, setOpen] = useState(false);
  const anchorRef = useRef<HTMLDivElement>(null);

  const classes = useStyles();

  const handleToggle = () => {
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event: any) => {
    if (
      anchorRef.current &&
      anchorRef.current.contains(event.target as HTMLElement)
    ) {
      return;
    }

    setOpen(false);
  };

  const [switchIdentitiesOpen, setSwitchIdentitiesOpen] =
    useState<boolean>(false);

  const { openSnackbar } = useSnackbar();

  const closeTab = () => {
    window.opener = null;
    window.open('', '_self');
    window.close();
  };

  const handleShutdown = () => {
    const sendShutdownSignal = async () => {
      return await CommonService.shutdown();
    };

    sendShutdownSignal()
      .then((response) => {
        closeTab();
      })
      .catch((err) => {
        openSnackbar('Unable to shutdown service gracefully', 'error');
      });
  };

  const renderUserInfo = () => {
    if (userIdentity != null) {
      return (
        <Fragment>
          <Avatar
            variant="rounded"
            src={userIdentity.picture}
            className={classes.identityPicture}
          >
            {userIdentity.name.charAt(0)}
          </Avatar>
          <Box display={'flex'} ml={1.5}>
            <Typography className={classes.selectedIdentity}>
              {userIdentity.name}
            </Typography>
          </Box>
        </Fragment>
      );
    } else {
      return (
        <Fragment>
          <Avatar variant="rounded" className={classes.identityPicture}>
            ?
          </Avatar>
          <Box display={'flex'} ml={1.5}>
            <Typography className={classes.selectedIdentity}>
              No identity
            </Typography>
          </Box>
        </Fragment>
      );
    }
  };

  return (
    <div ref={anchorRef} className={classes.buttonWrapper}>
      <Button onClick={handleToggle}>
        <Box className={classes.userMenuWrapper}>
          {renderUserInfo()}
          {open ? (
            <KeyboardArrowUpRoundedIcon />
          ) : (
            <KeyboardArrowDownRoundedIcon />
          )}
        </Box>
      </Button>
      <Popper
        open={open}
        anchorEl={anchorRef.current}
        role={undefined}
        transition
      >
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{
              transformOrigin:
                placement === 'bottom' ? 'center top' : 'center bottom'
            }}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList
                  style={{
                    background: 'white',
                    borderRadius: '5px',
                    boxShadow: 'none'
                  }}
                >
                  <MenuItem
                    className={classes.userMenuItem}
                    onClick={() => {
                      setSwitchIdentitiesOpen(true);
                      handleToggle();
                    }}
                  >
                    <Box display={'flex'}>
                      <SwapHorizRoundedIcon />
                      <Box ml={1}>Change identity</Box>
                    </Box>
                  </MenuItem>
                  <MenuItem
                    className={classes.userMenuItem}
                    onClick={() => {
                      handleShutdown();
                    }}
                  >
                    <Box display={'flex'}>
                      <PowerSettingsNewRoundedIcon />
                      <Box ml={1}>Shutdown</Box>
                    </Box>
                  </MenuItem>
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
      <SwapIdentities
        modalOpen={switchIdentitiesOpen}
        setModalOpen={setSwitchIdentitiesOpen}
      />
    </div>
  );
};

const useStyles = makeStyles((theme: Theme) => {
  return {
    buttonWrapper: {
      marginLeft: 'auto'
    },
    userMenuWrapper: {
      display: 'flex',
      alignItems: 'center'
    },
    selectedIdentity: {
      fontWeight: 600
    },
    identityPicture: {
      width: '45px',
      height: '45px',
      boxShadow: theme.palette.boxShadows.darker
    },
    userMenuItem: {
      fontFamily: 'Montserrat',
      fontWeight: 500,
      fontSize: '0.875rem',
      color: 'black'
    }
  };
});

export default UserMenu;

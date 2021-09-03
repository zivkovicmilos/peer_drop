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
import React, { FC, useContext, useRef, useState } from 'react';
import SessionContext from '../../../context/SessionContext';
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

  return (
    <div ref={anchorRef} className={classes.buttonWrapper}>
      <Button onClick={handleToggle}>
        <Box className={classes.userMenuWrapper}>
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
                  <MenuItem className={classes.userMenuItem}>
                    <Box display={'flex'}>
                      <SwapHorizRoundedIcon />
                      <Box ml={1}>Change identity</Box>
                    </Box>
                  </MenuItem>
                  <MenuItem className={classes.userMenuItem}>
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

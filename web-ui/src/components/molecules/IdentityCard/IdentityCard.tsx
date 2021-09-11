import {
  Avatar,
  Box,
  IconButton,
  Menu,
  MenuItem,
  Typography
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { MoreVertRounded } from '@material-ui/icons';
import TodayRoundedIcon from '@material-ui/icons/TodayRounded';
import VpnKeyRoundedIcon from '@material-ui/icons/VpnKeyRounded';
import clsx from 'clsx';
import { FC, useContext, useState } from 'react';
import { useHistory } from 'react-router-dom';
import SessionContext from '../../../context/SessionContext';
import { ReactComponent as CurrentIdentity } from '../../../shared/assets/icons/verified_black_24dp.svg';
import { ReactComponent as WorkspacesRoundedIcon } from '../../../shared/assets/icons/workspaces_black_24dp.svg';
import theme from '../../../theme/theme';
import {
  EIdentityCardMenuItem,
  IIdentityCardProps
} from './identityCard.types';

const IdentityCard: FC<IIdentityCardProps> = (props) => {
  const { picture, name, publicKeyID, numWorkspaces, creationDate } = props;

  const { userIdentity } = useContext(SessionContext);

  const classes = useStyles();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleClick = (event: any) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  interface IdentityCardMenuItem {
    type: EIdentityCardMenuItem;
    onClick: (identityId: string) => void;
  }

  const history = useHistory();

  const menuItems: IdentityCardMenuItem[] = [
    {
      type: EIdentityCardMenuItem.EDIT,
      onClick: (identityId: string) => {
        history.push(`identities/${identityId}/edit`);
      }
    },
    {
      type: EIdentityCardMenuItem.SHARE,
      onClick: (identityId: string) => {
        // TODO define
      }
    },
    {
      type: EIdentityCardMenuItem.BACKUP,
      onClick: (identityId: string) => {
        // TODO define
      }
    },
    {
      type: EIdentityCardMenuItem.SET_IDENTITY,
      onClick: (identityId: string) => {
        // TODO define
      }
    }
  ];

  return (
    <Box className={classes.identityCardWrapper}>
      <Box
        display={'flex'}
        justifyContent={'space-between'}
        alignItems={'center'}
      >
        <Box display={'flex'} alignItems={'center'} width={'100%'}>
          <Avatar src={picture}>{name.charAt(0)}</Avatar>
          <Box ml={1} width={'auto'} className={'truncate'} maxWidth={'140px'}>
            <Typography className={clsx(classes.identityCardName, 'truncate')}>
              {name}
            </Typography>
          </Box>

          {userIdentity.keyID == publicKeyID && (
            <Box ml={0.5}>
              <CurrentIdentity
                style={{
                  width: '15px',
                  height: 'auto'
                }}
              />
            </Box>
          )}
        </Box>
        <Box>
          <IconButton onClick={handleClick}>
            <MoreVertRounded
              style={{
                fill: 'black',
                width: '20px',
                height: 'auto'
              }}
            />
          </IconButton>
        </Box>
      </Box>
      <Box display={'flex'} flexDirection={'column'} width={'100%'} mt={2}>
        <Box
          display={'flex'}
          width={'100%'}
          className={'truncate'}
          alignItems={'center'}
          mb={0.5}
        >
          <VpnKeyRoundedIcon
            style={{
              fill: 'black',
              width: '18px',
              height: 'auto'
            }}
          />
          <Box ml={1} className={'truncate'} width={'100%'} maxWidth={'160px'}>
            <Typography className={clsx(classes.identitySubtext, 'truncate')}>
              {publicKeyID}
            </Typography>
          </Box>
        </Box>

        <Box
          display={'flex'}
          width={'100%'}
          className={'truncate'}
          alignItems={'center'}
          mb={0.5}
        >
          <WorkspacesRoundedIcon
            style={{
              fill: 'black',
              width: '18px',
              height: 'auto'
            }}
          />
          <Box ml={1} className={'truncate'} width={'100%'} maxWidth={'160px'}>
            <Typography className={clsx(classes.identitySubtext, 'truncate')}>
              {`${numWorkspaces} workspaces`}
            </Typography>
          </Box>
        </Box>

        <Box
          display={'flex'}
          width={'100%'}
          className={'truncate'}
          alignItems={'center'}
        >
          <TodayRoundedIcon
            style={{
              fill: 'black',
              width: '18px',
              height: 'auto'
            }}
          />
          <Box ml={1} className={'truncate'} width={'100%'} maxWidth={'160px'}>
            <Typography className={clsx(classes.identitySubtext, 'truncate')}>
              {creationDate}
            </Typography>
          </Box>
        </Box>
        <Menu
          anchorEl={anchorEl}
          open={Boolean(anchorEl)}
          onClose={handleClose}
        >
          {menuItems.map((menuItem) => {
            return (
              <MenuItem
                onClick={() => {
                  menuItem.onClick(publicKeyID);
                  handleClose();
                }}
                className={classes.identityMenuItem}
              >
                {menuItem.type}
              </MenuItem>
            );
          })}
        </Menu>
      </Box>
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    identityCardWrapper: {
      boxShadow: theme.palette.boxShadows.darker,
      display: 'flex',
      flexDirection: 'column',
      padding: '10px 15px',
      borderRadius: '15px',
      width: '250px',
      height: '160px',
      justifyContent: 'center',
      marginLeft: '36px',
      marginTop: '36px'
    },
    identityCardName: {
      fontWeight: 600
    },
    identitySubtext: {
      fontSize: '0.875rem'
    },
    identityMenuItem: {
      fontFamily: 'Montserrat',
      fontWeight: 500,
      fontSize: '0.875rem',
      color: 'black'
    }
  };
});

export default IdentityCard;

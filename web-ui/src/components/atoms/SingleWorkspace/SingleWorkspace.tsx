import { Box, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import React, { FC, useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import WorkspacesService from '../../../services/workspaces/workspacesService';
import greenCircle from '../../../shared/assets/img/workspace_green.png';
import redCircle from '../../../shared/assets/img/workspace_red.png';
import CommonUtils from '../../../shared/utils/CommonUtils';
import theme from '../../../theme/theme';
import useSnackbar from '../../molecules/Snackbar/useSnackbar.hook';
import { ISingleWorkspaceProps } from './singleWorkspace.types';

const SingleWorkspace: FC<ISingleWorkspaceProps> = (props) => {
  const { title, id, mnemonic } = props;

  const classes = useStyles();

  const generateBackgroundColor = () => {
    switch (getRandomNum(0, 5)) {
      case 0:
        return theme.palette.workspaceGradients.lightBlue;
      case 1:
        return theme.palette.workspaceGradients.lightGreen;
      case 2:
        return theme.palette.workspaceGradients.lightYellow;
      case 3:
        return theme.palette.workspaceGradients.lightPink;
      case 4:
        return theme.palette.workspaceGradients.lightPurple;
      case 5:
        return theme.palette.workspaceGradients.lightBrown;
    }
  };

  const [numPeersText, setNumPeersText] = useState<string>('');
  const [background, setBackground] = useState<string>('');

  enum CIRCLE_STATUS {
    GREEN,
    RED
  }

  const [numPeers, setNumPeers] = useState<number>(0);
  const [circle, setCircle] = useState<CIRCLE_STATUS>(CIRCLE_STATUS.RED);

  useEffect(() => {
    setBackground(generateBackgroundColor());
  }, []);

  useEffect(() => {
    setNumPeersText(renderNumPeers(numPeers));

    if (numPeers > 0) {
      setCircle(CIRCLE_STATUS.GREEN);
    } else {
      setCircle(CIRCLE_STATUS.RED);
    }
  }, [numPeers]);

  const { openSnackbar } = useSnackbar();

  const fetchNumPeers = () => {
    const fetchPeers = async () => {
      return await WorkspacesService.getWorkspacePeers(
        CommonUtils.formatMnemonic(mnemonic)
      );
    };

    fetchPeers()
      .then((response) => {
        setNumPeers(response.numPeers);
      })
      .catch((err) => {
        openSnackbar('Unable to fetch peer count', 'error');
      });
  };

  // Set the peer number updater
  useEffect(() => {
    fetchNumPeers();

    const peerNumUpdater = setInterval(() => {
      fetchNumPeers();
    }, 3000);

    return () => clearInterval(peerNumUpdater);
  }, []);

  const renderNumPeers = (numPeers: number) => {
    if (numPeers == 0) {
      return 'No peers';
    }

    if (numPeers < 2) {
      return '1 peer';
    } else {
      return `${numPeers} peers`;
    }
  };

  const getRandomNum = (min: number, max: number) => {
    min = Math.ceil(min);
    max = Math.floor(max);

    return Math.floor(Math.random() * (max - min + 1)) + min;
  };

  const history = useHistory();
  const handleWorkspaceClick = () => {
    history.push('/workspaces/view/' + CommonUtils.formatMnemonic(mnemonic));
  };

  return (
    <Box
      className={classes.workspaceItemWrapper}
      style={{
        background: background
      }}
      onClick={handleWorkspaceClick}
    >
      <Box ml={'auto'} display={'flex'} alignItems={'center'}>
        <Typography className={classes.workspaceInfo}>
          {numPeersText}
        </Typography>
        <Box ml={1} display={'flex'} alignItems={'center'}>
          {circle == CIRCLE_STATUS.GREEN ? (
            <img
              src={greenCircle}
              style={{
                width: '15px',
                height: 'auto'
              }}
            />
          ) : (
            <img
              src={redCircle}
              style={{
                width: '15px',
                height: 'auto'
              }}
            />
          )}
        </Box>
      </Box>
      <Box mt={'auto'}>
        <Typography className={classes.workspaceTitle}>{title}</Typography>
      </Box>
    </Box>
  );
};

const useStyles = makeStyles(() => {
  return {
    workspaceItemWrapper: {
      borderRadius: '15px',
      height: '135px',
      width: '210px',
      display: 'flex',
      padding: '10px 15px',
      marginLeft: '50px',
      marginTop: '50px',
      flexDirection: 'column',
      boxShadow: theme.palette.boxShadows.main,
      cursor: 'pointer',
      transition: 'box-shadow .24s ease-in-out',
      '&:hover': {
        boxShadow: theme.palette.boxShadows.darker
      }
    },
    workspaceTitle: {
      fontWeight: 600,
      textAlign: 'left'
    },
    workspaceInfo: {
      fontWeight: 400,
      fontSize: theme.typography.pxToRem(12)
    }
  };
});

export default SingleWorkspace;

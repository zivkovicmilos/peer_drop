import {
  Box,
  IconButton,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import DeleteRoundedIcon from '@material-ui/icons/DeleteRounded';
import clsx from 'clsx';
import React, { FC, useEffect, useState } from 'react';
import RendezvousService from '../../../services/rendezvous/rendezvousService';
import greenCircle from '../../../shared/assets/img/workspace_green.png';
import redCircle from '../../../shared/assets/img/workspace_red.png';
import theme from '../../../theme/theme';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import NoData from '../../atoms/NoData/NoData';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import {
  IRendezvousListProps,
  IRendezvousStatusResponse
} from './rendezvousList.types';

const RendezvousList: FC<IRendezvousListProps> = (props) => {
  const { openSnackbar } = useSnackbar();

  const { trigger, setTrigger } = props;

  const classes = useStyles();

  const [count, setCount] = useState<number>(0);
  const [rendezvous, setRendezvous] =
    useState<{ data: IRendezvousStatusResponse[] }>(null);

  const renderStatus = (status: string) => {
    if (status == 'on') {
      return (
        <Box display={'flex'} alignItems={'center'}>
          <Typography>Connected</Typography>
          <Box ml={2} display={'flex'} alignItems={'center'}>
            <img
              src={greenCircle}
              style={{
                width: '15px',
                height: 'auto'
              }}
            />
          </Box>
        </Box>
      );
    } else {
      return (
        <Box display={'flex'} alignItems={'center'}>
          <Typography>Disconnected</Typography>
          <Box ml={2} display={'flex'} alignItems={'center'}>
            <img
              src={redCircle}
              style={{
                width: '15px',
                height: 'auto'
              }}
            />
          </Box>
        </Box>
      );
    }
  };

  const fetchRendezvousNodes = () => {
    const fetchNodes = async () => {
      return await RendezvousService.getRendezvousNodes();
    };

    fetchNodes()
      .then((response) => {
        setCount(response.count);
        setRendezvous({ data: response.data });
      })
      .catch((err) => {
        openSnackbar('Unable to fetch rendezvous node list', 'error');
      });
  };

  // Set the connected status updater
  useEffect(() => {
    fetchRendezvousNodes();

    const rendezvousUpdater = setInterval(() => {
      fetchRendezvousNodes();
    }, 6000);

    return () => clearInterval(rendezvousUpdater);
  }, []);

  useEffect(() => {
    fetchRendezvousNodes();
  }, [trigger]);

  const handleDelete = (address: string) => {
    const deleteNode = async () => {
      return await RendezvousService.removeRendezvousNode({ address });
    };

    deleteNode()
      .then((response) => {
        setTrigger(!trigger);
        openSnackbar('Rendezvous successfully deleted', 'success');
      })
      .catch((err) => {
        openSnackbar('Unable to delete rendezvous node', 'error');
      });
  };

  const renderActions = (address: string) => {
    return (
      <Box display={'flex'} alignItems={'center'} justifyContent={'center'}>
        <IconButton
          style={{
            marginLeft: '1rem'
          }}
          classes={{
            root: 'iconButtonRoot'
          }}
          onClick={() => {
            handleDelete(address);
          }}
        >
          <DeleteRoundedIcon
            style={{
              fill: 'black'
            }}
          />
        </IconButton>
      </Box>
    );
  };

  if (rendezvous && count > 0) {
    return (
      <Box display={'flex'} width={'100%'} flexDirection={'column'}>
        <TableContainer
          component={Paper}
          classes={{
            root: classes.tableContainer
          }}
        >
          <Table className={classes.table}>
            <TableHead className={classes.tableHeadWrapper}>
              <TableRow>
                <TableCell
                  className={clsx(classes.tableHead, classes.noBorder)}
                  align="left"
                >
                  Address
                </TableCell>
                <TableCell
                  className={clsx(classes.tableHead, classes.noBorder)}
                  align="center"
                >
                  Status
                </TableCell>
                <TableCell
                  className={clsx(classes.tableHead, classes.noBorder)}
                  align="center"
                >
                  Actions
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {rendezvous.data.map((row, index) => (
                <TableRow key={`${row.address}-${index}`}>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="left"
                  >
                    {row.address}
                  </TableCell>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="left"
                  >
                    {renderStatus(row.status)}
                  </TableCell>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="left"
                  >
                    {renderActions(row.address)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Box>
    );
  } else if (count == 0) {
    return <NoData text={'No rendezvous nodes found'} />;
  } else {
    return <LoadingIndicator style={{ color: theme.palette.primary.main }} />;
  }
};

const useStyles = makeStyles(() => {
  return {
    table: {},
    tableHead: {
      fontWeight: 600
    },
    tableHeadWrapper: {
      borderBottom: '3px solid #F6F6F6'
    },
    tableContainer: {
      boxShadow: 'none !important'
    },
    noBorder: {
      border: 'none'
    },
    tableCell: {
      fontWeight: 500
    }
  };
});

export default RendezvousList;

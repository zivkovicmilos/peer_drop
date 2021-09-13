import {
  Box,
  IconButton,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { GridCellParams, GridColDef } from '@material-ui/data-grid';
import DeleteRoundedIcon from '@material-ui/icons/DeleteRounded';
import EditRoundedIcon from '@material-ui/icons/EditRounded';
import clsx from 'clsx';
import { FC, useEffect } from 'react';
import { useHistory } from 'react-router-dom';
import theme from '../../../theme/theme';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import NoData from '../../atoms/NoData/NoData';
import Pagination from '../../atoms/Pagination/Pagination';
import usePagination from '../../atoms/Pagination/pagination.hook';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import { IContactsTableProps } from './contactsTable.types';

const ContactsTable: FC<IContactsTableProps> = (props) => {
  const history = useHistory();

  const { handleDelete } = props;

  const handleEdit = (contactId: string | number) => {
    history.push('/contacts/' + contactId + '/edit');
  };

  const columns: GridColDef[] = [
    {
      field: 'name',
      headerName: 'Name',
      width: 200,
      cellClassName: 'muiGridTableCell'
    },
    {
      field: 'email',
      headerName: 'Email',
      width: 230,
      cellClassName: 'muiGridTableCell'
    },
    {
      field: 'publicKeyID',
      headerName: 'Public Key ID',
      width: 200,
      cellClassName: 'muiGridTableCell'
    },
    {
      field: 'dateAdded',
      headerName: 'Date Added',
      width: 160,
      cellClassName: 'muiGridTableCell'
    },
    {
      field: 'actions',
      headerName: 'Actions',
      type: 'string',
      width: 130,
      cellClassName: 'muiGridTableCell',
      renderCell: (params: GridCellParams) => {
        const id = params.id as string;
        const name = params.getValue(params.id, 'name') as string;
        const publicKey = params.getValue(params.id, 'publicKeyID') as string;

        return (
          <Box display={'flex'} alignItems={'center'}>
            <IconButton
              classes={{
                root: 'iconButtonRoot'
              }}
              onClick={() => handleEdit(id)}
            >
              <EditRoundedIcon
                style={{
                  fill: 'black'
                }}
              />
            </IconButton>

            <IconButton
              style={{
                marginLeft: '1rem'
              }}
              classes={{
                root: 'iconButtonRoot'
              }}
              onClick={() => {
                handleDelete({ id, name, publicKey });
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
      }
    }
  ];

  const rows = [
    {
      id: '1',
      name: 'Milos Zivkovic',
      email: 'milos@zmilos.com',
      publicKeyID: '4AEE18F83AFDEB23',
      dateAdded: '21.01.2020.'
    },
    {
      id: '2',
      name: 'Milos Zivkovic',
      email: 'milos@zmilos.com',
      publicKeyID: '4AEE18F83AFDEB23',
      dateAdded: '21.01.2020.'
    },
    {
      id: '3',
      name: 'Milos Zivkovic',
      email: 'milos@zmilos.com',
      publicKeyID: '4AEE18F83AFDEB23',
      dateAdded: '21.01.2020.'
    },
    {
      id: '4',
      name: 'Milos Zivkovic',
      email: 'milos@zmilos.com',
      publicKeyID: '4AEE18F83AFDEB23',
      dateAdded: '21.01.2020.'
    },
    {
      id: '5',
      name: 'Milos Zivkovic',
      email: 'milos@zmilos.com',
      publicKeyID: '4AEE18F83AFDEB23',
      dateAdded: '21.01.2020.'
    },
    {
      id: '6',
      name: 'Milos Zivkovic',
      email: 'milos@zmilos.com',
      publicKeyID: '4AEE18F83AFDEB23',
      dateAdded: '21.01.2020.'
    },
    {
      id: '7',
      name: 'Milos Zivkovic',
      email: 'milos@zmilos.com',
      publicKeyID: '4AEE18F83AFDEB23',
      dateAdded: '21.01.2020.'
    }
  ];

  const { page, count, setCount, limit, handlePageChange } = usePagination({
    limit: 8
  });

  const { openSnackbar } = useSnackbar();

  const classes = useStyles();

  const renderActions = (id: string, name: string, publicKey: string) => {
    return (
      <Box display={'flex'} alignItems={'center'} justifyContent={'center'}>
        <IconButton
          classes={{
            root: 'iconButtonRoot'
          }}
          onClick={() => handleEdit(id)}
        >
          <EditRoundedIcon
            style={{
              fill: 'black'
            }}
          />
        </IconButton>

        <IconButton
          style={{
            marginLeft: '1rem'
          }}
          classes={{
            root: 'iconButtonRoot'
          }}
          onClick={() => {
            handleDelete({ id, name, publicKey });
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

  useEffect(() => {
    setCount(rows.length);
  }, []);

  if (rows && count > 0) {
    return (
      <Box display={'flex'} width={'100%'} flexDirection={'column'}>
        <TableContainer
          component={Paper}
          classes={{
            root: classes.tableContainer
          }}
        >
          <Table className={classes.table} aria-label="simple table">
            <TableHead className={classes.tableHeadWrapper}>
              <TableRow>
                <TableCell
                  className={clsx(classes.tableHead, classes.noBorder)}
                  align="center"
                >
                  Name
                </TableCell>
                <TableCell
                  className={clsx(classes.tableHead, classes.noBorder)}
                  align="center"
                >
                  Email
                </TableCell>
                <TableCell
                  className={clsx(classes.tableHead, classes.noBorder)}
                  align="center"
                >
                  Public Key ID
                </TableCell>
                <TableCell
                  className={clsx(classes.tableHead, classes.noBorder)}
                  align="center"
                >
                  Date Added
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
              {rows.map((row, index) => (
                <TableRow key={`${row.name}-${index}`}>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="center"
                  >
                    {row.name}
                  </TableCell>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="center"
                  >
                    {row.email}
                  </TableCell>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="center"
                  >
                    {row.publicKeyID}
                  </TableCell>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="center"
                  >
                    {row.dateAdded}
                  </TableCell>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="center"
                  >
                    {renderActions(row.id, row.name, row.publicKeyID)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
        <Box
          width={'100%'}
          display={'flex'}
          alignItems={'center'}
          justifyContent={'center'}
        >
          <Pagination
            count={count}
            limit={limit}
            page={page}
            onPageChange={handlePageChange}
          />
        </Box>
      </Box>
    );
  } else if (count == 0) {
    return <NoData text={'No contacts found'} />;
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

export default ContactsTable;

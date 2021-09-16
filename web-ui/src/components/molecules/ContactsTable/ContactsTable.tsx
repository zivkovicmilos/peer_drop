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
import DeleteRoundedIcon from '@material-ui/icons/DeleteRounded';
import EditRoundedIcon from '@material-ui/icons/EditRounded';
import clsx from 'clsx';
import { FC, useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import ContactsService from '../../../services/contacts/contactsService';
import { IContactResponse } from '../../../services/contacts/contactsService.types';
import theme from '../../../theme/theme';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import NoData from '../../atoms/NoData/NoData';
import Pagination from '../../atoms/Pagination/Pagination';
import usePagination from '../../atoms/Pagination/pagination.hook';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import { IContactsTableProps } from './contactsTable.types';

const ContactsTable: FC<IContactsTableProps> = (props) => {
  const history = useHistory();

  const { handleDelete, fetchTrigger } = props;

  const handleEdit = (contactId: string | number) => {
    history.push('/contacts/' + contactId + '/edit');
  };

  const [contacts, setContacts] = useState<{ data: IContactResponse[] }>(null);

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
    const fetchContacts = async () => {
      return await ContactsService.getContacts({ page, limit });
    };

    fetchContacts()
      .then((contactsResponse) => {
        setCount(contactsResponse.count);
        setContacts({ data: contactsResponse.data });
      })
      .catch((err) => {
        setCount(0);

        openSnackbar('Unable to fetch contacts', 'error');
      });
  }, [page, fetchTrigger]);

  if (contacts && count > 0) {
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
              {contacts.data.map((row, index) => (
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

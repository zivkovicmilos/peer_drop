import { Box, IconButton } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { DataGrid, GridCellParams, GridColDef } from '@material-ui/data-grid';
import DeleteRoundedIcon from '@material-ui/icons/DeleteRounded';
import EditRoundedIcon from '@material-ui/icons/EditRounded';
import { FC, useState } from 'react';
import { useHistory } from 'react-router-dom';
import theme from '../../../theme/theme';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import NoData from '../../atoms/NoData/NoData';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import { IContactsTableProps } from './contactsTable.types';

const ContactsTable: FC<IContactsTableProps> = () => {
  const history = useHistory();

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
        const id = params.id;

        return (
          <Box display={'flex'} alignItems={'center'}>
            <IconButton onClick={() => handleEdit(id)}>
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

  const [count, setCount] = useState<number>(rows.length); // - 1
  const [page, setPage] = useState<number>(1); // 1

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };

  const { openSnackbar } = useSnackbar();

  if (rows && count > 0) {
    // todo change back
    return (
      <div style={{ width: '100%' }}>
        <DataGrid
          rowHeight={40}
          autoHeight={true}
          disableColumnFilter={true}
          disableColumnMenu={true}
          disableSelectionOnClick={true}
          rows={rows}
          columns={columns.map((column) => ({
            ...column,
            disableClickEventBubbling: true
          }))}
          isRowSelectable={() => false}
          pageSize={5}
          rowCount={count}
          paginationMode="server"
          onPageChange={handlePageChange}
        />
      </div>
    );
  } else if (count == 0) {
    return <NoData text={'No contacts found'} />;
  } else {
    return <LoadingIndicator style={{ color: theme.palette.primary.main }} />;
  }
};

const useStyles = makeStyles(() => {
  return {};
});

export default ContactsTable;

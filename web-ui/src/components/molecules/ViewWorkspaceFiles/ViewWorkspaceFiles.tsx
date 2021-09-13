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
import GetAppRoundedIcon from '@material-ui/icons/GetAppRounded';
import clsx from 'clsx';
import moment from 'moment';
import { FC, useEffect } from 'react';
import FileIcon, { IconStyle } from 'react-fileicons';
import folderIcon from '../../../shared/assets/img/folder.png';
import ColorUtils from '../../../shared/utils/ColorUtils';
import CommonUtils from '../../../shared/utils/CommonUtils';
import theme from '../../../theme/theme';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import NoData from '../../atoms/NoData/NoData';
import Pagination from '../../atoms/Pagination/Pagination';
import usePagination from '../../atoms/Pagination/pagination.hook';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import {
  FILE_TYPE,
  IViewWorkspaceFilesProps
} from './viewWorkspaceFiles.types';

const ViewWorkspaceFiles: FC<IViewWorkspaceFilesProps> = (props) => {
  const { workspaceInfo } = props;

  const dummyDate = new Date();

  const files = [
    {
      id: '1',
      name: 'PoS Roadmap',
      type: FILE_TYPE.FILE,
      extension: 'doc',
      modified: dummyDate,
      size: 123123
    },
    {
      id: '2',
      name: 'IBFT Draft',
      type: FILE_TYPE.FILE,
      extension: 'pdf',
      modified: dummyDate,
      size: 123123
    },
    {
      id: '3',
      name: 'Onboarding',
      type: FILE_TYPE.FOLDER,
      modified: dummyDate,
      size: 123123
    },
    {
      id: '4',
      name: 'Tasks',
      type: FILE_TYPE.FILE,
      extension: 'zip',
      modified: dummyDate,
      size: 123123
    },
    {
      id: '5',
      name: '.api-keys',
      type: FILE_TYPE.FILE,
      extension: '.api-keys',
      modified: dummyDate,
      size: 123123
    }
  ];

  const { page, count, setCount, limit, handlePageChange } = usePagination({
    limit: 8
  });

  const { openSnackbar } = useSnackbar();

  const classes = useStyles();

  const handleDownload = (id: string) => {};

  const renderActions = (id: string) => {
    return (
      <Box display={'flex'} alignItems={'center'} justifyContent={'center'}>
        <IconButton
          classes={{
            root: 'iconButtonRoot'
          }}
          onClick={() => {
            handleDownload(id);
          }}
        >
          <GetAppRoundedIcon
            style={{
              fill: 'black'
            }}
          />
        </IconButton>
      </Box>
    );
  };

  useEffect(() => {
    setCount(files.length);
  }, []);

  const renderIcon = (extension: string) => {
    if (!extension) {
      return (
        <Box mr={2}>
          <img src={folderIcon} className={classes.folderIcon} />
        </Box>
      );
    }

    const colorCode = ColorUtils.getColorCode(extension);

    return (
      <Box mr={2}>
        <FileIcon
          extension={extension}
          background={'white'}
          colorScheme={{ primary: colorCode.iconColor }}
          iconStyle={IconStyle.normal}
          size={35}
        />
      </Box>
    );
  };

  const renderItemName = (
    type: FILE_TYPE,
    name: string,
    extension?: string
  ) => {
    if (type == FILE_TYPE.FILE) {
      return `${name}.${extension}`;
    }

    return name;
  };

  if (files && count > 0) {
    return (
      <Box display={'flex'} width={'80%'} flexDirection={'column'}>
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
                  Modified
                </TableCell>
                <TableCell
                  className={clsx(classes.tableHead, classes.noBorder)}
                  align="center"
                >
                  Size
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
              {files.map((row, index) => (
                <TableRow key={`${row.name}-${index}`}>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="center"
                  >
                    <Box display={'flex'} alignItems={'center'}>
                      {renderIcon(row.extension)}
                      <Typography className={classes.itemName}>
                        {renderItemName(row.type, row.name, row.extension)}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="center"
                  >
                    {moment(row.modified).format('DD.MM.YYYY.')}
                  </TableCell>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="center"
                  >
                    {CommonUtils.formatBytes(row.size)}
                  </TableCell>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="center"
                  >
                    {renderActions(row.id)}
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
    itemName: {
      fontSize: '0.875rem',
      fontWeight: 500
    },
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
      borderBottom: '3px solid #F6F6F6'
    },
    folderIcon: {
      width: '35px',
      height: 'auto'
    },
    tableCell: {
      fontWeight: 500
    }
  };
});

export default ViewWorkspaceFiles;

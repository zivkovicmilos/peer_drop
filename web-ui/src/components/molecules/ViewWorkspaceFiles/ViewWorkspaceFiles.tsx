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
import { FC, useEffect, useState } from 'react';
import FileIcon, { IconStyle } from 'react-fileicons';
import { IWorkspaceDetailedFileResponse } from '../../../services/workspaces/workspacesService.types';
import folderIcon from '../../../shared/assets/img/folder.png';
import ColorUtils from '../../../shared/utils/ColorUtils';
import CommonUtils from '../../../shared/utils/CommonUtils';
import theme from '../../../theme/theme';
import LoadingIndicator from '../../atoms/LoadingIndicator/LoadingIndicator';
import NoData from '../../atoms/NoData/NoData';
import Pagination from '../../atoms/Pagination/Pagination';
import usePagination from '../../atoms/Pagination/pagination.hook';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import { IViewWorkspaceFilesProps } from './viewWorkspaceFiles.types';

const ViewWorkspaceFiles: FC<IViewWorkspaceFilesProps> = (props) => {
  const { workspaceInfo } = props;

  // Pagination is done locally
  const [shownFiles, setShownFiles] = useState<{
    files: IWorkspaceDetailedFileResponse[];
  }>(null);

  const { page, count, setCount, limit, handlePageChange } = usePagination({
    limit: 8
  });

  useEffect(() => {
    if (workspaceInfo.workspaceFiles.length > 0) {
      setShownFiles({ files: workspaceInfo.workspaceFiles.slice(0, limit) });
    }

    setCount(workspaceInfo.workspaceFiles.length);
  }, [workspaceInfo]);

  useEffect(() => {
    if (workspaceInfo.workspaceFiles.length > 0) {
      let offset = (page - 1) * limit;

      let upperBound = offset + limit;
      if (upperBound > workspaceInfo.workspaceFiles.length) {
        upperBound = workspaceInfo.workspaceFiles.length;
      }

      setShownFiles({
        files: workspaceInfo.workspaceFiles.slice(offset, upperBound)
      });
    }
  }, [page]);

  const { openSnackbar } = useSnackbar();

  const classes = useStyles();

  const handleDownload = (checksum: string) => {
  };

  const renderActions = (checksum: string) => {
    return (
      <Box display={'flex'} alignItems={'center'} justifyContent={'center'}>
        <IconButton
          classes={{
            root: 'iconButtonRoot'
          }}
          onClick={() => {
            handleDownload(checksum);
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

  const renderItemName = (name: string, extension: string) => {
    return `${name}.${extension}`;
  };

  if (shownFiles && count > 0) {
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
              {shownFiles.files.map((row, index) => (
                <TableRow key={`${row.name}-${index}`}>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="center"
                  >
                    <Box display={'flex'} alignItems={'center'}>
                      {renderIcon(row.extension)}
                      <Typography className={classes.itemName}>
                        {renderItemName(row.name, row.extension)}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell
                    className={clsx(classes.noBorder, classes.tableCell)}
                    align="center"
                  >
                    {moment.unix(row.dateModified).format('DD.MM.YYYY.')}
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
                    {renderActions(row.checksum)}
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
    return <NoData text={'No files found'} />;
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

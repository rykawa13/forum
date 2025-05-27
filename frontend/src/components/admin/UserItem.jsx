
import React from 'react';
import { ListItem, ListItemAvatar, Avatar, ListItemText, ListItemSecondaryAction, IconButton, Chip } from '@mui/material';
import { Delete as DeleteIcon, Edit as EditIcon } from '@mui/icons-material';

const UserItem = ({ user, onDelete, onEdit }) => {
  return (
    <ListItem>
      <ListItemAvatar>
        <Avatar>{user.username.charAt(0)}</Avatar>
      </ListItemAvatar>
      <ListItemText
        primary={user.username}
        secondary={user.email}
      />
      {user.isAdmin && (
        <Chip 
          label="Админ" 
          color="primary" 
          size="small" 
          sx={{ mr: 2 }}
        />
      )}
      <ListItemSecondaryAction>
        <IconButton edge="end" onClick={() => onEdit(user.id)} sx={{ mr: 1 }}>
          <EditIcon />
        </IconButton>
        <IconButton edge="end" onClick={() => onDelete(user.id)} color="error">
          <DeleteIcon />
        </IconButton>
      </ListItemSecondaryAction>
    </ListItem>
  );
};

export default UserItem;
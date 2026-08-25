import {
  Container,
  TextField,
  Typography,
  type TextFieldProps,
} from "@mui/material";

type XTextFieldProps = TextFieldProps & {
  labelName: string;
};

export default function XTextField({ labelName, ...props }: XTextFieldProps) {
  return (
    <Container>
      <label htmlFor={props.id}>
        <Typography variant="h6">{labelName}</Typography>
      </label>
      <TextField {...props} variant="outlined" fullWidth color="primary" />
    </Container>
  );
}

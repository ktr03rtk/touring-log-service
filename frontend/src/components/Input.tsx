type InputProperties = {
  name: string;
  title: string;
  type: string;
  value?: string | number;
  handleChange: (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => void;
  placeholder?: string;
  className?: string;
  errorDiv?: string;
  errorMsg?: string;
};

const Input = (props: InputProperties) => {
  return (
    <div className='mb-3'>
      <label htmlFor={props.name} className='form-label'>
        {props.title}
      </label>
      <input
        type={props.type}
        className={`form-control ${props.className}`}
        id={props.name}
        name={props.name}
        value={props.value}
        onChange={props.handleChange}
        placeholder={props.placeholder}
      />
      <div className={props.errorDiv}>{props.errorMsg}</div>
    </div>
  );
};

export default Input;

import { InputHTMLAttributes, useId } from 'react'

type InputProps = InputHTMLAttributes<HTMLInputElement> & {
  label?: string
  inputSize?: 'md' | 'lg' | 'sm'
  wrapperClass?: string
  inputClass?: string
}

const Input = (props: InputProps) => {
  const { label, inputSize, inputClass, wrapperClass: cln, ...rest } = props
  const id = useId()

  const _class = cln || 'mb-6 h-full'
  const sizeClass =
    inputSize === 'sm' ? 'py-1 px-2' : inputSize === 'lg' ? 'py-4 px-6' : 'py-3 px-4'

  return (
    <div className={_class}>
      {label && <label htmlFor={id}>{label}</label>}
      <input
        id={id}
        {...rest}
        className={`w-full text-sm border overflow-hidden rounded border-gray-500 outline-cyan-500 ${sizeClass} ${
          inputClass || ''
        }`}
      />
    </div>
  )
}

export default Input

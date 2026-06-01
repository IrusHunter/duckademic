import { BiPlus, BiX } from "react-icons/bi";
import { useState } from 'react';
import type { FieldDef } from '../../types/admin';
import { RelationSelect } from '../RelationSelect/RelationSelect';
import styles from './AddForm.module.css';

type Props = {
  fields: FieldDef[]
  onSubmit: (data: Record<string, string>) => void
}

export function AddForm({ fields, onSubmit }: Props) {
  const [values, setValues] = useState<Record<string, string>>({})
  const [open, setOpen] = useState(false)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(values)
    setValues({})
    setOpen(false)
  }

  return (
    <div className={styles.wrapper}>
      <button
        onClick={() => setOpen(!open)}
        className={styles.toggleButton}
      >
        {open ? (
          <><BiX size={20} /> Cancel</>
        ) : (
          <><BiPlus size={20} /> Add</>
        )}
      </button>
      {open && (
        <form onSubmit={handleSubmit} className={styles.form}>
          {fields.map(field => (
            <div key={field.key} className={styles.field}>
              <label className={styles.label}>{field.label}{field.required && ' *'}</label>
              {field.relation ? (
                <RelationSelect
                  field={field}
                  value={values[field.key] ?? ''}
                  onChange={v => setValues(prev => ({ ...prev, [field.key]: v }))}
                />
              ) : (
                <input
                  value={values[field.key] ?? ''}
                  onChange={e => setValues(prev => ({ ...prev, [field.key]: e.target.value }))}
                  required={field.required}
                  className={styles.input}
                />
              )}
            </div>
          ))}
          <div className={styles.actions}>
            <button type="submit" className={styles.submitButton}>
              Save
            </button>
          </div>
        </form>
      )}
    </div>
  )
}
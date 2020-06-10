import React from 'react';

import { Form } from '../components/Form';
import { Heading } from '../components/Heading';
import { Input } from '../components/Input';
import { Link } from '../components/Link';
import { List } from '../components/List';
import { fetch } from '../util/api';

{{~
  func javascriptValueify
    case $0
      when "integer", "bigint", "smallint", "decimal", "numeric", "real", "double precision"
        "Number"
      when "boolean"
        "Boolean"
      else
        ""
    end
  end
~}}

export function Create{{ table.name|string.capitalize }}() {
  const [state, setState] = React.useState({
    {{~ for column in table.columns ~}}
    {{~ if column.auto_increment
          continue
        end ~}}
    '{{ column.name }}': '',
    {{~ end ~}}
  });

  const [error, setError] = React.useState('');
  const handleSubmit = React.useCallback(async (e) => {
    e.preventDefault();
    setError('');

    try {
      const rsp = await fetch('{{ table.name }}', {
        {{~ for column in table.columns ~}}
        {{~ if column.auto_increment
              continue
            end ~}}
        '{{ column.name }}': {{ javascriptValueify column }}(state['{{ column.name }}']),
        {{~ end ~}}
      });

      if (rsp.error) {
        setError(rsp.error);
        return;
      }

      window.location = '/{{ table.name }}';
    } catch (e) {
      // Need the try-catch so we can return false here.
      console.error(e);
      return false;
    }
  });

  return (
    <>
      <Link to="/{{ table.name }}">{{ table.name|string.capitalize }}</Link>
      <Heading size="xl">Create</Heading>
      <Form error={error} buttonText="Create" onSubmit={handleSubmit}>
        {{~ for column in table.columns ~}}
        {{~ if column.auto_increment
              continue
            end ~}}
        <div className="mb-4">
          <Input
            label="{{ column.name }}"
            id="{{ column.name }}"
            value={state['{{ column.name }}']}
            onChange={(e) => {
              // e.target.value is not available within the setState callback, so copy it.
              // https://duncanleung.com/fixing-react-warning-synthetic-events-in-setstate/
              const { value } = e.target;
              setState(s => ({ ...s, ['{{ column.name }}']: value }))
            }}
          />
        </div>
        {{ end }}
      </Form>
    </>
  );
}

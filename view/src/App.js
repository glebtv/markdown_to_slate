import React, { Component } from 'react';

import prettyBytes from 'pretty-bytes';
import { Value } from 'slate'
import { Editor } from 'slate-react';
import data from "./data.json"
import './App.css';

console.log("data:")
console.log(data)

const initialValue = {
  "document": {
    "object": "document",
    "data": {},
    "nodes": data
  }
}

class Image extends React.PureComponent {
  render() {
    const { node } = this.props
    const { data } = node
    const src = data.get("src")
    const { attributes } = this.props
    return src ? (
      <img alt="" src={src} ref={(r) => this.image = r} />
    ) : (
      <div {...attributes}>Loading...</div>
    )
  }
}

class App extends Component {
  state = {
    value: Value.fromJSON(initialValue)
  };

  onChange = ({ value }) => {
    this.setState({ value });
  };

  renderNode(props, editor, next) {
    const { attributes, children, node, isFocused } = props
    const { data } = node

    //console.log("render node", node.type, node)

    switch (node.type) {
      case 'code':
        return <code {...props}>
          {children}
        </code>
      case 'code_line':
        return <div {...props}>
          {children}
        </div>

      case 'check-list-item':
        return (
          <div>
            {children}
          </div>
        )

      case 'link':
        const href = data.get('href')
        //console.log("link", children)
        return (
          <a {...attributes} href={href} target="_blank" rel="nofollow noopener noreferrer">
            {children}
          </a>
        )

      case 'mention':
        const userID = data.get('user')
        return (
          <a
            contentEditable={false}
            {...attributes}
          >
            {userID}
          </a>
        )


      case 'issue':
        const issueID = data.get('issue')
        return (
          <a
            selected={isFocused}
            contentEditable={false}
            {...attributes}
          >
            задача {issueID}
            {children}
          </a>
        )

      case 'file':
        return (
          <a
            contentEditable={false}
            href={data.get('url')}
            {...attributes}
          >
            {data.get('name')} ({prettyBytes(data.get("size"))})
            {children}
          </a>
        )

      case 'image':
        return <Image selected={isFocused} {...props} />
      case 'block-quote':
        return <blockquote {...attributes}>{children}</blockquote>
      case 'bulleted-list':
        return <ul {...attributes}>{children}</ul>
      case 'heading-one':
        return <h1 {...attributes}>{children}</h1>
      case 'heading-two':
        return <h2 {...attributes}>{children}</h2>
      case 'heading-three':
        return <h3 {...attributes}>{children}</h3>
      case 'list-item':
        return <li {...attributes}>{children}</li>
      case 'numbered-list':
        return <ol {...attributes}>{children}</ol>
      case 'div':
        return <div {...attributes}>{children}</div>
      case 'paragraph':
        return <div {...attributes}>{children}</div>
      default:
        return next()
    }
  }

  /**
   * Render a Slate mark.
   *
   * @param {Object} props
   * @return {Element}
   */
  renderMark(props, editor, next) {
    const { children, mark, attributes } = props

    switch (mark.type) {
       // Code
       case 'comment':
          return (
            <span {...attributes} style={{ opacity: '0.33' }}>
              {children}
            </span>
          )
        case 'keyword':
          return (
            <span {...attributes} style={{ fontWeight: 'bold' }}>
              {children}
            </span>
          )
        case 'tag':
          return (
            <span {...attributes} style={{ fontWeight: 'bold' }}>
              {children}
            </span>
          )
        case 'punctuation':
          return (
            <span {...attributes} style={{ opacity: '0.75' }}>
              {children}
            </span>
          )
      // End code

      case 'bold':
        return <strong {...attributes}>{children}</strong>
      case 'code':
        return <code {...attributes}>{children}</code>
      case 'italic':
        return <em {...attributes}>{children}</em>
      case 'underlined':
        return <u {...attributes}>{children}</u>
      default:
        return next()
    }
  }


  render() {
    //console.log('render', this.state);

    return (
      <div className="App">
        <Editor
          value={this.state.value}
          onChange={this.onChange}
          renderNode={this.renderNode}
          renderMark={this.renderMark}
        />
      </div>
    );
  }
}

export default App;

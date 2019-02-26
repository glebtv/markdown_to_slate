import React, { Component } from 'react';

import prettyBytes from 'pretty-bytes';
import { Value } from 'slate'
import { Editor } from 'slate-react';
//import data from "./data.json"
import './App.css';

//console.log("data:")
//console.log(data)


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

function slate2value(slate) {
  return Value.fromJSON({
    "document": {
      "object": "document",
      "data": {},
      "nodes": slate
    }
  })
}

class SlatePreview extends Component {
  state = {
    value: slate2value(this.props.value)
  };

  onChange = ({ value }) => {
    this.setState({ value });
  };

  componentWillReceiveProps(props) {
    this.setState({value: slate2value(props.value)})
  }

  renderNode(props, editor, next) {
    const { attributes, children, node, isFocused } = props
    const { data } = node

    //console.log("render node", node.type, node)

    switch (node.type) {
      case 'code':
        return <code {...attributes}>
          {children}
        </code>
      case 'code_line':
        return <div {...attributes}>
          {children}
        </div>

      case 'check-list-item':
        const checked = data.get("checked")
        return (
          <div className="check-list-item" {...attributes}>
            <input type="checkbox" checked={checked} onChange={() => {}}  {...attributes} />
            {children}
          </div>
        )

      case 'link':
        const href = data.get('href')
        //console.log("link", href, children)
        return (
          <a {...attributes} href={href} target="_blank" rel="nofollow noopener noreferrer">
            {children}
          </a>
        )

      case 'mention':
        const userID = data.get('user')
        console.log("mention")
        return (
          <a
            contentEditable={false}
            {...attributes}
          >
            mention @{userID}
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
        //console.log("render div", children)
        //return <div {...attributes}>div: {children} /div</div>
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
        <button type="button" onClick={(e) => {
          console.log(this.state.value.toJSON().document.nodes)
        }}>log state</button>
      </div>
    );
  }
}

export default SlatePreview;

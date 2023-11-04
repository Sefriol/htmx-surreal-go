package view

var Index = `<!DOCTYPE html>
<html lang="en">
  <head>
    <title></title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.4" integrity="sha384-zUfuhFKKZCbHTY6aRR46gxiqszMk5tcHjsVFxnUo8VMus4kHGVdIYVbOYYNlKmHV" crossorigin="anonymous"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@1/css/pico.min.css">
    <link href="css/style.css" rel="stylesheet">
    <script src="dist/bundle.js"></script>
    </head>
    <body>
      <main class="grid">
        <div class="container">
          {{template "Users" . }}
        </div>
        <div class="container"> 
	  <orb-graph></orb-graph>
	  <div id="relations"></div>
        </div>
    </main>
  </body>
</html>`
var Users = `{{define "Users"}}
<form hx-post="/user" hx-target="#users" hx-swap="afterbegin">
  <label>
    Name
        <input name="name" type="text">  
  </label>
  <label>
    Surname 
        <input name="surname" type="text">  
  </label>
  <button type="submit">Create</button>
</form>
<div id="users">
{{range .Users}}
    {{template "user" .}}
{{end}}
</div>
{{end}}`

var Relations = `{{define "relations"}}
<div>
{{ if . }}{{.}}
<h1>Relatives for {{(index . 0).User.Name}} {{(index . 0).User.Surname}}</h1>
{{else}}
	<h1>User has no relatives</h1>
{{end}}
<table>
  <thead>
    <tr>
      <th scope="col">#</th>
      <th scope="col">Child</th>
      <th scope="col">Parent</th>
      <th scope="col">Actions</th>
    </tr>
  </thead>
  <tbody id="relation-table">
  {{range $index,$element:= .}}
    <tr class="relation-row">
      <th scope="row">{{$index}}</th>
      <td>{{$element.Child.ID}}</td>
      <td>{{$element.Parent.ID}}</td>
      <td><a hx-trigger="click" hx-delete="/user/{{$element.Child.ID}}/relative/{{$element.Parent.ID}}" hx-target="closest .relation-row" hx-confirm="Are you sure you wish to delete this relation?">Delete</a></td>
    </tr>
  {{end}}
  </tbody>
</table>
</div>
{{end}}`

var Relation = `{{define "relation"}}
<tr>
  <th scope="row">{{.Index}}</th>
  <td>{{.Element.Child.ID}}</td>
  <td>{{.Element.Parent.ID}}</td>
</tr>
{{end}}`

var EditUser = `{{define "edit-user"}}
<form hx-put="/user/{{.ID}}" hx-target="this" hx-swap="outerHTML">
  <div>
    <label>First Name</label>
    <input type="text" name="name" value="{{.Name}}">
  </div>
  <div class="form-group">
    <label>Last Name</label>
    <input type="text" name="surname" value="{{.Surname}}">
  </div>
  <button class="btn">Submit</button>
  <button class="btn" hx-get="/user/{{.ID}}">Cancel</button>
</form>
{{end}}`


var User = `{{define "user"}}
<div class="user" hx-target="#relations" hx-trigger="click" hx-get="/user/{{.ID}}/relatives">
  <article>
    <div class="id"> ID: {{.ID}}</div>
    <div class="name"> Name: {{.Name}}</div>
    <div class="surname"> Surname: {{.Surname}}</div>
    <a hx-trigger="click" hx-delete="/user/{{.ID}}" hx-target="closest .user" hx-confirm="Are you sure you wish to delete your account?">Delete</a>
    <a hx-trigger="click" hx-get="/user/{{.ID}}/edit" hx-target="closest .user">Edit</a>
    <a hx-trigger="click" hx-get="/user/{{.ID}}/relative" hx-target="closest .user" hx-swap="afterend">Add Relative</a>
  </article>
</div>
{{end}}
`

var RelationDialog = `{{define "relation-dialog"}}
<dialog open>
  <main>
    <section>
      <form hx-post="/user/{{.User.ID}}/relative" hx-target="#relations">
        <select name="relation">
          <option value="child">Child</option>
          <option value="parent">Parent</option>
        </select>
        <select name="relative">
        {{range .Users}}
          <option value="{{.ID}}">{{.Name}} {{.Surname}}</option>
        {{end}}
        </select>
        <button type="submit">Submit</button>
      </form>
    </section>
    <section>
      <form>
        <button class="secondary" formmethod="dialog">Cancel</button>
      </form>
    </section>
  </main>
</dialog>
{{end}}
`

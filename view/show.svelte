<script>
  export let blog = {}
</script>

<h1>Show Blog</h1>

<table border="1" cellpadding="10">
  <thead>
    {#each Object.keys(blog) as key}
      <th>{key}</th>
    {/each}
  </thead>
  <tr>
    {#each Object.keys(blog) as key}
      <td>{blog[key]}</td>
    {/each}
  </tr>
</table>

<br />

<a href={`/`}>Back</a>
<span> | </span>
<a href={`/${blog.id || 0}/edit`}>Edit</a>

<style>
  table {
    border-collapse: collapse;
  }
</style>

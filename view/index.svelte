<script>
  export let blogs = []
</script>

<h1>Blog Index</h1>

<table border="1" cellpadding="10">
  {#if blogs.length > 0}
    <thead>
      {#each Object.keys(blogs[0]) as key}
        <th>{key}</th>
      {/each}
    </thead>
  {/if}
  {#each blogs as blog}
    <tr>
      {#each Object.keys(blog) as key}
        {#if key.toLowerCase() === "id"}
          <td><a href={`/${blog.id || 0}`}>{blog[key]}</a></td>
        {:else}
          <td>{blog[key]}</td>
        {/if}
      {/each}
    </tr>
  {/each}
</table>

<br />

<a href={`/new`}>New Blog</a>

<style>
  table {
    border-collapse: collapse;
  }
</style>

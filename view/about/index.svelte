<script>
  import Head from "../Head.svelte";
  import Foot from "../Foot.svelte";
</script>

<Head />

<body>
  <div class="containter">
    <h1 class="display-1 text-center">🤔 about me</h1>

    <div class="containter-md">
      <div class="row align-items-center">
        <div class="col text-center">
          <img src="../../public/images/headshot.jpg" alt="me" />
        </div>
        <div class="col">
          <ul>
            <li>
              Hi I'm Warren! I'm a backend software engineer based in the south
              east of England.
            </li>
            <li>
              I graduated the University of Brighton in 2022 and have been
              working since then as a backend developer in golang, working on
              microservice architectures deployed on kubernetes and
              containterized with docker.
            </li>
            <li>
              The purpose of this website is to act as a dumping ground for
              things that I've learnt in my last 2 years as well as things I'm
              currenltly learning about or things that I'm struggling with as an
              backend engineer.
            </li>
            <li>
              I'm planning to try and update the website every 2 weeks and will
              add a log of the updates made so you can follow along 😄 Thanks
              for stopping by and I hope you enjoy!
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</body>

<Foot />

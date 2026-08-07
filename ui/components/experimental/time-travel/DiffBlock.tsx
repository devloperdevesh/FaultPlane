"use client";


interface Props {

type:
"added"
|
"removed"
|
"modified";

label:string;

before?:string;

after?:string;

}



const styles = {

added:
"border-emerald-500/20 bg-emerald-500/10",

removed:
"border-red-500/20 bg-red-500/10",

modified:
"border-yellow-500/20 bg-yellow-500/10",

};



export default function DiffBlock({

type,

label,

before,

after,

}:Props){


return (

<div

className={`
rounded-xl
border
p-4
${styles[type]}
`}

>


<div
className="
flex
justify-between
"
>

<p
className="
text-sm
font-medium
text-white
"
>

{label}

</p>


<span
className="
text-xs
uppercase
text-zinc-400
"
>

{type}

</span>


</div>



<div
className="
mt-4
space-y-2
font-mono
text-xs
"
>


{
before &&

<p
className="
text-red-400
"
>
- {before}

</p>

}



{
after &&

<p
className="
text-emerald-400
"
>
+ {after}

</p>

}



</div>


</div>

)

}